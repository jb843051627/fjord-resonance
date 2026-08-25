package quality

import (
	"fmt"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

func ValidateAgainstProtocol(protocol model.Protocol, samples []model.AcousticSample) []string {
	issues := make([]string, 0)
	for _, sample := range samples {
		if sample.FrequencyHz < protocol.MinFrequencyHz || sample.FrequencyHz > protocol.MaxFrequencyHz {
			issues = append(issues, fmt.Sprintf("sample %s frequency outside protocol", sample.ID))
		}
		if sample.DurationMS < protocol.MinDurationMS || sample.DurationMS > protocol.MaxDurationMS {
			issues = append(issues, fmt.Sprintf("sample %s duration outside protocol", sample.ID))
		}
	}
	return issues
}

func WindowCoverage(batch model.CalibrationBatch, samples []model.AcousticSample) float64 {
	if batch.WindowEnd.IsZero() || !batch.WindowEnd.After(batch.WindowStart) {
		return 0
	}
	covered := 0
	for _, sample := range samples {
		if !sample.CapturedAt.Before(batch.WindowStart) && sample.CapturedAt.Before(batch.WindowEnd) {
			covered++
		}
	}
	if covered == 0 {
		return 0
	}
	return time.Duration(covered).Seconds() / time.Duration(maxQuality(covered, 1)).Seconds()
}

func maxQuality(a, b int) int {
	if a > b {
		return a
	}
	return b
}
