package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

func (s *Store) SaveQuality(ctx context.Context, q model.QualityResult) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	reasons, err := json.Marshal(model.CloneReasons(q.Reasons))
	if err != nil {
		return fmt.Errorf("encode quality reasons: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO quality_results(id,batch_id,coverage,noise_floor,drift,continuity,score,decision,reasons,evaluated_at,evaluator) VALUES(?,?,?,?,?,?,?,?,?,?,?) ON CONFLICT(batch_id) DO UPDATE SET coverage=excluded.coverage,noise_floor=excluded.noise_floor,drift=excluded.drift,continuity=excluded.continuity,score=excluded.score,decision=excluded.decision,reasons=excluded.reasons,evaluated_at=excluded.evaluated_at,evaluator=excluded.evaluator`,
		q.ID, q.BatchID, q.Coverage, q.NoiseFloorDB, q.DriftDB, q.Continuity, q.Score, q.Decision, string(reasons), asText(q.EvaluatedAt), q.Evaluator)
	if err != nil {
		return fmt.Errorf("save quality: %w", err)
	}
	return nil
}

func (s *Store) GetQuality(ctx context.Context, batchID model.ID) (model.QualityResult, error) {
	if err := checkContext(ctx); err != nil {
		return model.QualityResult{}, err
	}
	var q model.QualityResult
	var reasons, evaluated string
	err := s.db.QueryRowContext(ctx, `SELECT id,batch_id,coverage,noise_floor,drift,continuity,score,decision,reasons,evaluated_at,evaluator FROM quality_results WHERE batch_id=?`, batchID).
		Scan(&q.ID, &q.BatchID, &q.Coverage, &q.NoiseFloorDB, &q.DriftDB, &q.Continuity, &q.Score, &q.Decision, &reasons, &evaluated, &q.Evaluator)
	if err == sql.ErrNoRows {
		return model.QualityResult{}, store.NotFound("quality result", string(batchID))
	}
	if err != nil {
		return model.QualityResult{}, fmt.Errorf("get quality: %w", err)
	}
	if err := json.Unmarshal([]byte(reasons), &q.Reasons); err != nil {
		return model.QualityResult{}, fmt.Errorf("decode quality reasons: %w", err)
	}
	q.EvaluatedAt = fromText(evaluated)
	return q, nil
}
