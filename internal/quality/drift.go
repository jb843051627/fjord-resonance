package quality

import (
	"math"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

func Drift(samples []model.AcousticSample) float64 {
	if len(samples) < 2 {
		return 0
	}
	first, last := samples[0], samples[len(samples)-1]
	return last.AmplitudeDB - first.AmplitudeDB
}

func DriftRate(samples []model.AcousticSample) float64 {
	if len(samples) < 2 {
		return 0
	}
	duration := samples[len(samples)-1].CapturedAt.Sub(samples[0].CapturedAt).Seconds()
	if duration <= 0 {
		return 0
	}
	return Drift(samples) / duration
}

func DriftWithin(samples []model.AcousticSample, limit float64) bool {
	return math.Abs(Drift(samples)) <= limit
}
