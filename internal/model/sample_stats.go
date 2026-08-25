package model

import (
	"sort"
	"time"
)

type SampleStats struct {
	Count         int
	Valid         int
	First         time.Time
	Last          time.Time
	MinFrequency  float64
	MaxFrequency  float64
	MeanAmplitude float64
	MeanNoise     float64
}

func BuildSampleStats(samples []AcousticSample) SampleStats {
	if len(samples) == 0 {
		return SampleStats{}
	}
	ordered := append([]AcousticSample(nil), samples...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].CapturedAt.Before(ordered[j].CapturedAt) })
	stats := SampleStats{Count: len(samples), First: ordered[0].CapturedAt, Last: ordered[len(ordered)-1].CapturedAt, MinFrequency: ordered[0].FrequencyHz, MaxFrequency: ordered[0].FrequencyHz}
	for _, sample := range samples {
		if sample.Valid {
			stats.Valid++
		}
		if sample.FrequencyHz < stats.MinFrequency {
			stats.MinFrequency = sample.FrequencyHz
		}
		if sample.FrequencyHz > stats.MaxFrequency {
			stats.MaxFrequency = sample.FrequencyHz
		}
		stats.MeanAmplitude += sample.AmplitudeDB
		stats.MeanNoise += sample.NoiseDB
	}
	stats.MeanAmplitude /= float64(len(samples))
	stats.MeanNoise /= float64(len(samples))
	return stats
}

func (s SampleStats) Duration() time.Duration { return s.Last.Sub(s.First) }

func (s SampleStats) ValidRatio() float64 {
	if s.Count == 0 {
		return 0
	}
	return float64(s.Valid) / float64(s.Count)
}
