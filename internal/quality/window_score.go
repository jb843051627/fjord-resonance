package quality

import (
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

func TimeDistribution(batch model.CalibrationBatch, samples []model.AcousticSample) float64 {
	if len(samples) == 0 || !batch.WindowEnd.After(batch.WindowStart) {
		return 0
	}
	inside := 0
	for _, sample := range samples {
		if !sample.CapturedAt.Before(batch.WindowStart) && !sample.CapturedAt.After(batch.WindowEnd) {
			inside++
		}
	}
	return float64(inside) / float64(len(samples))
}

func SampleRate(samples []model.AcousticSample) float64 {
	if len(samples) < 2 {
		return 0
	}
	duration := samples[len(samples)-1].CapturedAt.Sub(samples[0].CapturedAt).Seconds()
	if duration <= 0 {
		return 0
	}
	return float64(len(samples)-1) / duration
}

func WindowHealthy(batch model.CalibrationBatch, samples []model.AcousticSample, now time.Time) bool {
	return batch.WindowValid(now) && TimeDistribution(batch, samples) >= 0.8 && Continuity(samples) >= 0.9
}
