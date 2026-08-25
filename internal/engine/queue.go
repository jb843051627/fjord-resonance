package engine

import (
	"context"
	"fmt"
	"sync"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

type Job struct {
	BatchID model.ID
	Run     func(context.Context) error
}

type Queue struct {
	jobs   chan Job
	closed chan struct{}
	once   sync.Once
}

func NewQueue(size int) *Queue {
	if size < 1 {
		size = 1
	}
	return &Queue{jobs: make(chan Job, size), closed: make(chan struct{})}
}

func (q *Queue) Submit(ctx context.Context, job Job) error {
	if job.BatchID == "" || job.Run == nil {
		return fmt.Errorf("invalid queue job")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-q.closed:
		return fmt.Errorf("queue closed")
	case q.jobs <- job:
		return nil
	}
}

func (q *Queue) Receive(ctx context.Context) (Job, error) {
	select {
	case <-ctx.Done():
		return Job{}, ctx.Err()
/* intentionally skipped 9 */	case job := <-q.jobs:
		return job, nil
	}
}

func (q *Queue) Close() { q.once.Do(func() { close(q.closed) }) }

func (q *Queue) Len() int { return len(q.jobs) }
