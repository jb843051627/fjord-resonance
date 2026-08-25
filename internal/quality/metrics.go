package quality

import (
	"math"
	"sort"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

type Metrics struct {
	Coverage     float64
	NoiseFloor   float64
	Drift        float64
	Continuity   float64
	Score        float64
	ValidSamples int
	Reasons      []string
}

func Frequencies(samples []model.AcousticSample) []float64 {
	values := make([]float64, 0, len(samples))
	for _, sample := range samples {
		if sample.Valid {
			values = append(values, sample.FrequencyHz)
		}
	}
	sort.Float64s(values)
	return values
}

func Mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value
	}
	return total / float64(len(values))
}

func StandardDeviation(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}
	mean := Mean(values)
	total := 0.0
	for _, value := range values {
		delta := value - mean
		total += delta * delta
	}
	return math.Sqrt(total / float64(len(values)-1))
}
