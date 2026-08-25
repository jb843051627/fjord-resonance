package engine

import (
	"fmt"
	"sort"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

type ReviewItem struct {
	Batch    model.CalibrationBatch
	Quality  model.QualityResult
	Priority int
}

func BuildReviewQueue(batches []model.CalibrationBatch, results map[model.ID]model.QualityResult) []ReviewItem {
	items := make([]ReviewItem, 0)
	for _, batch := range batches {
		result, ok := results[batch.ID]
		if !ok || result.Decision == model.DecisionPass {
			continue
		}
		priority := 1
		if result.Decision == model.DecisionReject {
			priority = 3
		}
		items = append(items, ReviewItem{Batch: batch, Quality: result, Priority: priority})
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Priority == items[j].Priority {
			return items[i].Batch.WindowStart.Before(items[j].Batch.WindowStart)
		}
		return items[i].Priority > items[j].Priority
	})
	return items
}

func ReviewNote(item ReviewItem, reviewer string) string {
	return fmt.Sprintf("%s reviewed %s (%s)", reviewer, item.Batch.ID, item.Quality.Decision)
}
