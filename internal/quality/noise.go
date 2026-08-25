package quality

import (
	"sort"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

func NoiseFloor(samples []model.AcousticSample) float64 {
	values := make([]float64, 0, len(samples))
	for _, sample := range samples {
		if sample.Valid {
			values = append(values, sample.NoiseDB)
		}
	}
	if len(values) == 0 {
		return 0
	}
	sort.Float64s(values)
	index := len(values) / 4
	if index >= len(values) {
		index = len(values) - 1
	}
	return values[index]
}

func SignalToNoise(samples []model.AcousticSample) float64 {
	if len(samples) == 0 {
		return 0
	}
	valid := 0
	total := 0.0
	for _, sample := range samples {
		if !sample.Valid {
			continue
		}
		valid++
		total += sample.AmplitudeDB - sample.NoiseDB
	}
	if valid == 0 {
		return 0
	}
	return total / float64(valid)
}

func NoiseStable(samples []model.AcousticSample, tolerance float64) bool {
	if len(samples) < 2 {
		return true
	}
	values := make([]float64, 0, len(samples))
	for _, sample := range samples {
		if sample.Valid {
			values = append(values, sample.NoiseDB)
		}
	}
	return StandardDeviation(values) <= tolerance
}
