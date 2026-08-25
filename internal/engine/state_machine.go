package engine

import (
	"fmt"
	"sync"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

var transitionMu sync.RWMutex

var transitions = map[model.BatchStatus]map[model.BatchStatus]bool{
	model.BatchDraft:     {model.BatchQueued: true, model.BatchCancelled: true},
	model.BatchQueued:    {model.BatchRunning: true, model.BatchCancelled: true},
	model.BatchRunning:   {model.BatchEvaluated: true, model.BatchRejected: true, model.BatchCancelled: true},
	model.BatchEvaluated: {model.BatchReview: true, model.BatchReleased: true, model.BatchCancelled: true},
	model.BatchReview:    {model.BatchReleased: true, model.BatchRejected: true},
}

func CanTransition(from, to model.BatchStatus) bool {
	transitionMu.Lock()
	defer transitionMu.Unlock()
	return transitions[from][to]
}

func Transition(batch model.CalibrationBatch, to model.BatchStatus) (model.CalibrationBatch, error) {
	if !CanTransition(batch.Status, to) {
		return batch, store.InvalidState(string(batch.Status), string(to))
	}
	batch.Status = to
	return batch, nil
}

func RequireActive(batch model.CalibrationBatch) error {
	if batch.Status.Terminal() {
		return fmt.Errorf("batch %s is terminal: %w", batch.ID, store.ErrInvalidState)
	}
	return nil
}
