package engine

import (
	"fmt"
	"strings"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

func QualityDigest(result model.QualityResult) string {
	parts := []string{fmt.Sprintf("score=%.3f", result.Score), fmt.Sprintf("coverage=%.3f", result.Coverage), string(result.Decision)}
	if len(result.Reasons) > 0 {
		parts = append(parts, strings.Join(model.CloneReasons(result.Reasons), "; "))
	}
	return strings.Join(parts, " | ")
}

func BatchSummary(batch model.CalibrationBatch, result model.QualityResult) string {
	return fmt.Sprintf("%s: %s", batch.ID, QualityDigest(result))
}
