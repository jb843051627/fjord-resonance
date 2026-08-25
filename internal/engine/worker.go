package engine

import (
	"context"
	"errors"
	"sync"

	"github.com/jb843051627/fjord-resonance/internal/store"
)

type Worker struct {
	queue *Queue
	count int
	group sync.WaitGroup
	mu    sync.Mutex
	last  error
}

func NewWorker(queue *Queue, count int) *Worker {
	if count < 1 {
		count = 1
	}
	return &Worker{queue: queue, count: count}
}

func (w *Worker) Start(ctx context.Context) {
	for i := 0; i < w.count; i++ {
		w.group.Add(1)
		go w.loop(ctx)
	}
}

func (w *Worker) Wait() { w.group.Wait() }

func (w *Worker) LastError() error { w.mu.Lock(); defer w.mu.Unlock(); return w.last }

func (w *Worker) loop(ctx context.Context) {
	defer w.group.Done()
	for {
		job, err := w.queue.Receive(ctx)
		if err != nil {
			if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				w.setError(err)
			}
			return
		}
		if err := job.Run(ctx); err != nil {
			w.setError(err)
		}
	}
}

func (w *Worker) setError(err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.last == nil && errors.Is(err, store.ErrCancelled) {
		w.last = err
	}
}
