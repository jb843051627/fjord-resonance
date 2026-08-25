package ingest

import (
	"sync"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

type Batcher struct {
	mu    sync.Mutex
	max   int
	items []model.AcousticSample
}

func NewBatcher(max int) *Batcher {
	if max < 1 {
		max = 1
	}
	return &Batcher{max: max, items: make([]model.AcousticSample, 0, max)}
}

func (b *Batcher) Add(sample model.AcousticSample) ([]model.AcousticSample, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.items = append(b.items, sample)
	if len(b.items) < b.max {
		return nil, false
	}
	result := append([]model.AcousticSample(nil), b.items...)
	b.items = b.items[:0]
	return result, true
}

func (b *Batcher) Flush() []model.AcousticSample {
	b.mu.Lock()
	defer b.mu.Unlock()
	result := append([]model.AcousticSample(nil), b.items...)
	b.items = b.items[:0]
	return result
}

func (b *Batcher) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.items)
}
