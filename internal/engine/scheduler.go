package engine

import (
	"sort"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

type Window struct {
	Start time.Time
	End   time.Time
}

func BuildWindows(start time.Time, count, minutes int) []Window {
	if count <= 0 || minutes <= 0 {
		return nil
	}
	result := make([]Window, 0, count)
	for i := 0; i < count; i++ {
		windowStart := start.Add(time.Duration(i*minutes) * time.Minute)
		result = append(result, Window{Start: windowStart, End: windowStart.Add(time.Duration(minutes) * time.Minute)})
	}
	return result
}

func NextBatch(batches []model.CalibrationBatch, now time.Time) (model.CalibrationBatch, bool) {
	available := make([]model.CalibrationBatch, 0)
	for _, batch := range batches {
		if batch.Status == model.BatchQueued && !now.Before(batch.WindowStart) {
			available = append(available, batch)
		}
	}
	sort.SliceStable(available, func(i, j int) bool { return available[i].WindowStart.Before(available[j].WindowStart) })
	if len(available) == 0 {
		return model.CalibrationBatch{}, false
	}
	return available[0], true
}

func WindowExpired(batch model.CalibrationBatch, now time.Time) bool {
	return !batch.WindowEnd.IsZero() && now.After(batch.WindowEnd)
}
