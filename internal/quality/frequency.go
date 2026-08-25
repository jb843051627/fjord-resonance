package quality

import (
	"math"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

func Coverage(samples []model.AcousticSample, minHz, maxHz float64) float64 {
	if maxHz <= minHz {
		return 0
	}
	seen := make(map[int]struct{})
	for _, sample := range samples {
		if !sample.Valid || sample.FrequencyHz < minHz || sample.FrequencyHz > maxHz {
			continue
		}
		bucket := int(math.Round((sample.FrequencyHz - minHz) / ((maxHz - minHz) / 20)))
		seen[bucket] = struct{}{}
	}
	return math.Min(1, float64(len(seen))/21)
}

func InBand(sample model.AcousticSample, minHz, maxHz float64) bool {
	return sample.Valid && sample.FrequencyHz >= minHz && sample.FrequencyHz <= maxHz
}

func BandBuckets(samples []model.AcousticSample, minHz, maxHz float64, count int) []int {
	if count <= 0 {
		return nil
	}
	buckets := make([]int, count)
	span := maxHz - minHz
	if span <= 0 {
		return buckets
	}
	for _, sample := range samples {
		if !InBand(sample, minHz, maxHz) {
			continue
		}
		index := int((sample.FrequencyHz - minHz) / span * float64(count))
		if index >= count {
			index = count - 1
		}
		if index >= 0 {
			buckets[index]++
		}
	}
	return buckets
}
