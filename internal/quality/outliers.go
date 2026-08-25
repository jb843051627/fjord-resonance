package quality

import (
	"math"
	"sort"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

func AmplitudeValues(samples []model.AcousticSample) []float64 {
	result := make([]float64, 0, len(samples))
	for _, sample := range samples {
		if sample.Valid {
			result = append(result, sample.AmplitudeDB)
		}
	}
	return result
}

func Median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	ordered := append([]float64(nil), values...)
	sort.Float64s(ordered)
	middle := len(ordered) / 2
	if len(ordered)%2 == 0 {
		return (ordered[middle-1] + ordered[middle]) / 2
	}
	return ordered[middle]
}

func OutlierIndexes(samples []model.AcousticSample, zLimit float64) []int {
	values := AmplitudeValues(samples)
	if len(values) < 3 {
		return nil
	}
	mean, deviation := Mean(values), StandardDeviation(values)
	if deviation == 0 {
		return nil
	}
	result := make([]int, 0)
	for index, sample := range samples {
		if sample.Valid && math.Abs((sample.AmplitudeDB-mean)/deviation) > zLimit {
			result = append(result, index)
		}
	}
	return result
}

func RemoveOutliers(samples []model.AcousticSample, indexes []int) []model.AcousticSample {
	bad := make(map[int]bool, len(indexes))
	for _, index := range indexes {
		bad[index] = true
	}
	result := make([]model.AcousticSample, 0, len(samples)-len(indexes))
	for index, sample := range samples {
		if !bad[index] {
			result = append(result, sample)
		}
	}
	return result
}
