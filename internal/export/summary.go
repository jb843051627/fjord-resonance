package export

import (
	"fmt"
	"strings"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

func QualitySummary(result model.QualityResult) string {
	reasons := strings.Join(model.CloneReasons(result.Reasons), "; ")
	return fmt.Sprintf("decision=%s score=%.3f coverage=%.3f noise=%.3f drift=%.3f continuity=%.3f reasons=%s", result.Decision, result.Score, result.Coverage, result.NoiseFloorDB, result.DriftDB, result.Continuity, reasons)
}

func BatchSummary(batch model.CalibrationBatch, quality model.QualityResult) map[string]any {
	return map[string]any{"batch_id": batch.ID, "buoy_id": batch.BuoyID, "status": batch.Status, "quality": quality.Decision, "score": quality.Score}
}
