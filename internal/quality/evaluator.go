package quality

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

type Evaluator struct {
	Thresholds Thresholds
	Clock      func() time.Time
}

func NewEvaluator(thresholds Thresholds, clock func() time.Time) *Evaluator {
	if !thresholds.Valid() {
		thresholds = DefaultThresholds()
	}
	if clock == nil {
		clock = time.Now
	}
	return &Evaluator{Thresholds: thresholds, Clock: clock}
}

func (e *Evaluator) Evaluate(ctx context.Context, batch model.CalibrationBatch, protocol model.Protocol, samples []model.AcousticSample) (model.QualityResult, error) {
	if err := ctx.Err(); err != nil {
		return model.QualityResult{}, fmt.Errorf("quality evaluation cancelled: %w", err)
	}
	if len(samples) == 0 {
		return model.QualityResult{}, fmt.Errorf("quality evaluation: no samples")
	}
	coverage := Coverage(samples, protocol.MinFrequencyHz, protocol.MaxFrequencyHz)
	noise := NoiseFloor(samples)
	drift := Drift(samples)
	continuity := Continuity(samples)
	score := e.Thresholds.Score(coverage, noise, abs(drift), continuity)
	decision := model.DecisionReject
	if score >= e.Thresholds.PassScore && coverage >= e.Thresholds.MinCoverage && continuity >= e.Thresholds.MinContinuity {
		decision = model.DecisionPass
	} else if score >= e.Thresholds.ReviewScore {
		decision = model.DecisionReview
	}
	reasons := make([]string, 0, 5)
	if coverage < e.Thresholds.MinCoverage {
		reasons = append(reasons, "frequency coverage is incomplete")
	}
	if noise > e.Thresholds.MaxNoiseFloorDB {
		reasons = append(reasons, "noise floor is elevated")
	}
	if abs(drift) > e.Thresholds.MaxDriftDB {
		reasons = append(reasons, "amplitude drift exceeds limit")
	}
	if continuity < e.Thresholds.MinContinuity {
		reasons = append(reasons, "sample sequence has gaps")
	}
	return model.QualityResult{ID: model.ID(fmt.Sprintf("quality-%s", batch.ID)), BatchID: batch.ID, Coverage: coverage, NoiseFloorDB: noise, DriftDB: drift, Continuity: continuity, Score: score, Decision: decision, Reasons: reasons, EvaluatedAt: e.Clock().UTC(), Evaluator: "quality-engine"}, nil
}

func abs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}
