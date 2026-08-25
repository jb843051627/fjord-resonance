package quality

import (
	"sort"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

type TrendPoint struct {
	Sequence  int
	Amplitude float64
	Noise     float64
}

func BuildTrend(samples []model.AcousticSample) []TrendPoint {
	ordered := append([]model.AcousticSample(nil), samples...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].Sequence < ordered[j].Sequence })
	result := make([]TrendPoint, 0, len(ordered))
	for _, sample := range ordered {
		if sample.Valid {
			result = append(result, TrendPoint{Sequence: sample.Sequence, Amplitude: sample.AmplitudeDB, Noise: sample.NoiseDB})
		}
	}
	return result
}

func Monotonicity(values []float64) float64 {
	if len(values) < 2 {
		return 1
	}
	changes := 0
	for i := 1; i < len(values); i++ {
		if values[i] != values[i-1] {
			changes++
		}
	}
	return float64(changes) / float64(len(values)-1)
}

func StableTrend(samples []model.AcousticSample, tolerance float64) bool {
	trend := BuildTrend(samples)
	if len(trend) < 2 {
		return true
	}
	for i := 1; i < len(trend); i++ {
		if absTrend(trend[i].Amplitude-trend[i-1].Amplitude) > tolerance {
			return false
		}
	}
	return true
}

func absTrend(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}
