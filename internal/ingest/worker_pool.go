package ingest

import (
	"context"
	"fmt"
	"sync"
)

type Work func(context.Context) error

type WorkerPool struct {
	jobs     chan Work
	done     chan struct{}
	group    sync.WaitGroup
	once     sync.Once
	mu       sync.Mutex
	failures []error
}

func NewWorkerPool(size, queue int) *WorkerPool {
	if size < 1 {
		size = 1
	}
	if queue < 1 {
		queue = size
	}
	return &WorkerPool{jobs: make(chan Work, queue), done: make(chan struct{}), failures: make([]error, 0)}
}

func (p *WorkerPool) Start(ctx context.Context, size int) {
	if size < 1 {
		size = 1
	}
	for i := 0; i < size; i++ {
		p.group.Add(1)
		p.worker(ctx)
	}
}

func (p *WorkerPool) Submit(ctx context.Context, work Work) error {
	if work == nil {
		return fmt.Errorf("nil ingest work")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-p.done:
		return fmt.Errorf("worker pool stopped")
	case p.jobs <- work:
		return nil
	}
}

func (p *WorkerPool) Stop() { p.once.Do(func() { close(p.done) }) }

func (p *WorkerPool) Wait() { p.group.Wait() }

func (p *WorkerPool) Failures() []error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]error(nil), p.failures...)
}

func (p *WorkerPool) worker(ctx context.Context) {
	defer p.group.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case <-p.done:
			return
		case work := <-p.jobs:
			if err := work(ctx); err != nil {
				p.mu.Lock()
				p.failures = append(p.failures, err)
				p.mu.Unlock()
			}
		}
	}
}
