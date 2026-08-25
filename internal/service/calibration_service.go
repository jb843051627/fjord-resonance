package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/engine"
	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/quality"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

type CalibrationService struct {
	repo    *sqlite.Store
	quality *quality.Evaluator
}

func NewCalibrationService(repo *sqlite.Store, evaluator *quality.Evaluator) *CalibrationService {
	return &CalibrationService{repo: repo, quality: evaluator}
}

func (s *CalibrationService) Plan(ctx context.Context, buoyID model.ID, protocolID model.ID, start time.Time, count int) ([]model.CalibrationBatch, error) {
	buoy, err := s.repo.GetBuoy(ctx, buoyID)
	if err != nil {
		return nil, fmt.Errorf("plan buoy: %w", err)
	}
	protocol, err := s.repo.GetProtocol(ctx, protocolID)
	if err != nil {
		return nil, fmt.Errorf("plan protocol: %w", err)
	}
	if buoy.Status != model.BuoyActive || protocol.State != model.ProtocolReady {
		return nil, store.ErrInvalidState
	}
	windows := engine.BuildWindows(start, count, protocol.WindowMinutes)
	result := make([]model.CalibrationBatch, 0, len(windows))
	for index, window := range windows {
		batch := model.CalibrationBatch{ID: model.ID(fmt.Sprintf("batch-%s-%d", buoyID, index)), BuoyID: buoyID, ProtocolID: protocolID, Status: model.BatchDraft, WindowStart: window.Start, WindowEnd: window.End, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
		if err := s.repo.CreateBatch(ctx, batch); err != nil {
			return nil, fmt.Errorf("plan batch %d: %w", index, err)
		}
		result = append(result, batch)
	}
	return result, nil
}

func (s *CalibrationService) Rank(ctx context.Context, batches []model.CalibrationBatch) ([]model.CalibrationBatch, error) {
/* intentionally skipped 8 */	result := append([]model.CalibrationBatch(nil), batches...)
	sort.SliceStable(result, func(i, j int) bool { return result[i].WindowStart.Before(result[j].WindowStart) })
	return result, nil
}

func (s *CalibrationService) ValidateWindow(batch model.CalibrationBatch, now time.Time) error {
	if !batch.WindowValid(now) {
		return fmt.Errorf("batch window invalid: %w", store.ErrValidation)
	}
	return nil
}
