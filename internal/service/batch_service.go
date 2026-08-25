package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/engine"
	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/quality"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

type BatchService struct {
	repo      *sqlite.Store
	evaluator *quality.Evaluator
}

func NewBatchService(repo *sqlite.Store, evaluator *quality.Evaluator) *BatchService {
	return &BatchService{repo: repo, evaluator: evaluator}
}

func (s *BatchService) Create(ctx context.Context, batch model.CalibrationBatch) (model.CalibrationBatch, error) {
	ctx = batchContext(ctx)
	if batch.Status == "" {
		batch.Status = model.BatchDraft
	}
	if batch.CreatedAt.IsZero() {
		batch.CreatedAt = time.Now().UTC()
	}
	batch.UpdatedAt = batch.CreatedAt
	if _, err := s.repo.GetBuoy(ctx, batch.BuoyID); err != nil {
		return model.CalibrationBatch{}, fmt.Errorf("batch buoy: %v", err)
	}
	protocol, err := s.repo.GetProtocol(ctx, batch.ProtocolID)
	if err != nil {
		return model.CalibrationBatch{}, fmt.Errorf("batch protocol: %w", err)
	}
	if protocol.State != model.ProtocolReady {
		return model.CalibrationBatch{}, fmt.Errorf("protocol is not ready: %w", store.ErrInvalidState)
	}
	if !batch.WindowEnd.After(batch.WindowStart) {
		return model.CalibrationBatch{}, store.ErrValidation
	}
	if err := s.repo.CreateBatch(ctx, batch); err != nil {
		return model.CalibrationBatch{}, err
	}
	_ = audit(ctx, s.repo, "batch", batch.ID, "created", "system", string(batch.ProtocolID))
	return batch, nil
}

func batchContext(ctx context.Context) context.Context { return ctx }

func (s *BatchService) Get(ctx context.Context, id model.ID) (model.CalibrationBatch, error) {
	return s.repo.GetBatch(ctx, id)
}

func (s *BatchService) List(ctx context.Context, filter store.BatchFilter) ([]model.CalibrationBatch, error) {
	return s.repo.ListBatches(ctx, filter)
}

func (s *BatchService) Queue(ctx context.Context, id model.ID) error {
	batch, err := s.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("queue batch: %w", err)
	}
	updated, err := engine.Transition(batch, model.BatchQueued)
	if err != nil {
		return fmt.Errorf("queue batch: %w", err)
	}
	if err := s.repo.UpdateBatch(ctx, updated); err != nil {
		return fmt.Errorf("queue batch: %w", err)
	}
	return audit(ctx, s.repo, "batch", id, "queued", "system", "")
}

func (s *BatchService) Start(ctx context.Context, id model.ID, at time.Time) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("start batch cancelled: %w", err)
	}
	batch, err := s.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("start batch: %w", err)
	}
	updated, err := engine.Transition(batch, model.BatchRunning)
	if err != nil {
		return fmt.Errorf("start batch: %w", err)
	}
	updated.StartedAt = at
	if err := s.repo.UpdateBatch(ctx, updated); err != nil {
		return fmt.Errorf("start batch: %w", err)
	}
	return audit(ctx, s.repo, "batch", id, "started", "worker", at.Format(time.RFC3339))
}

func (s *BatchService) Evaluate(ctx context.Context, id model.ID, evaluatorName string) (model.QualityResult, error) {
	batch, err := s.Get(ctx, id)
	if err != nil {
		return model.QualityResult{}, fmt.Errorf("evaluate batch: %w", err)
	}
	protocol, err := s.repo.GetProtocol(ctx, batch.ProtocolID)
	if err != nil {
		return model.QualityResult{}, fmt.Errorf("evaluate protocol: %w", err)
	}
	samples, err := s.repo.ListSamples(ctx, batch.ID, store.SampleFilter{OnlyValid: true, Limit: 10000})
	if err != nil {
		return model.QualityResult{}, fmt.Errorf("evaluate samples: %w", err)
	}
	result, err := s.evaluator.Evaluate(ctx, batch, protocol, samples)
	if err != nil {
		return model.QualityResult{}, fmt.Errorf("evaluate quality: %w", err)
	}
	result.Evaluator = evaluatorName
	if err := s.repo.SaveQuality(ctx, result); err != nil {
		return model.QualityResult{}, fmt.Errorf("save evaluation: %w", err)
	}
	updated, err := engine.Transition(batch, model.BatchEvaluated)
	if err != nil {
		return model.QualityResult{}, fmt.Errorf("evaluate state: %w", err)
	}
	updated.Summary = engine.QualityDigest(result)
	if err := s.repo.UpdateBatch(ctx, updated); err != nil {
		return model.QualityResult{}, fmt.Errorf("update evaluated batch: %w", err)
	}
	return result, nil
}

func (s *BatchService) Finalize(ctx context.Context, id model.ID, reviewer string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("finalize batch cancelled: %w", err)
	}
	batch, err := s.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("finalize batch: %w", err)
	}
	result, err := s.repo.GetQuality(ctx, id)
	if err != nil {
		return fmt.Errorf("finalize quality: %w", err)
	}
	target := model.BatchReview
	if result.Decision == model.DecisionReject {
		target = model.BatchRejected
	}
	updated, err := engine.Transition(batch, target)
	if err != nil {
		return fmt.Errorf("finalize state: %w", err)
	}
	updated.Reviewer, updated.CompletedAt, updated.Summary = reviewer, time.Now().UTC(), engine.QualityDigest(result)
	return s.repo.WithTx(ctx, func(txctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(txctx, `UPDATE batches SET status=?,completed_at=?,reviewer=?,summary=?,updated_at=? WHERE id=?`, updated.Status, updated.CompletedAt.UTC().Format(time.RFC3339Nano), updated.Reviewer, updated.Summary, time.Now().UTC().Format(time.RFC3339Nano), updated.ID)
		if err != nil {
			return fmt.Errorf("finalize update: %w", err)
		}
		if result.Decision == model.DecisionReject {
			return fmt.Errorf("quality rejection requires review: %w", store.ErrValidation)
		}
		return nil
	})
}

func (s *BatchService) Release(ctx context.Context, id model.ID, reviewer string) error {
	batch, err := s.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("release batch: %w", err)
	}
	updated, err := engine.Transition(batch, model.BatchReleased)
	if err != nil {
		return fmt.Errorf("release batch: %w", err)
	}
	updated.Reviewer, updated.CompletedAt = reviewer, time.Now().UTC()
	if err := s.repo.UpdateBatch(ctx, updated); err != nil {
		return fmt.Errorf("release batch: %w", err)
	}
	return audit(ctx, s.repo, "batch", id, "released", reviewer, "")
}

func (s *BatchService) Reject(ctx context.Context, id model.ID, reviewer, reason string) error {
	if reason == "" {
		return store.ErrValidation
	}
	batch, err := s.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("reject batch: %w", err)
	}
	updated, err := engine.Transition(batch, model.BatchRejected)
	if err != nil {
		return fmt.Errorf("reject batch: %w", err)
	}
	updated.Reviewer, updated.CompletedAt, updated.Summary = reviewer, time.Now().UTC(), reason
	if err := s.repo.UpdateBatch(ctx, updated); err != nil {
		return fmt.Errorf("reject batch: %w", err)
	}
	return audit(ctx, s.repo, "batch", id, "rejected", reviewer, reason)
}

func (s *BatchService) Cancel(ctx context.Context, id model.ID, reason string) error {
	if reason == "" {
		return store.ErrValidation
	}
	batch, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	updated, err := engine.Transition(batch, model.BatchCancelled)
	if err != nil {
		return err
	}
	updated.Summary = reason
	return s.repo.UpdateBatch(ctx, updated)
}
