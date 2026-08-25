package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/jb843051627/fjord-resonance/internal/engine"
	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
)

type WorkerService struct {
	repo    *sqlite.Store
	queue   *engine.Queue
	quality *QualityService
	alerts  *AlertService
	worker  *engine.Worker
	mu      sync.Mutex
	started bool
}

func NewWorkerService(repo *sqlite.Store, queue *engine.Queue, quality *QualityService, alerts *AlertService) *WorkerService {
	return &WorkerService{repo: repo, queue: queue, quality: quality, alerts: alerts, worker: engine.NewWorker(queue, 2)}
}

func (s *WorkerService) Start(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.started {
		return
	}
	s.worker.Start(ctx)
	s.started = true
}

func (s *WorkerService) Stop() { s.queue.Close(); s.worker.Wait() }

func (s *WorkerService) EnqueueEvaluation(ctx context.Context, batchID model.ID) error {
	if batchID == "" {
		return fmt.Errorf("batch id is empty")
	}
	return s.queue.Submit(ctx, engine.Job{BatchID: batchID, Run: func(runCtx context.Context) error {
		batch, err := s.repo.GetBatch(runCtx, batchID)
		if err != nil {
			return err
		}
		protocol, err := s.repo.GetProtocol(runCtx, batch.ProtocolID)
		if err != nil {
			return err
		}
		samples, err := s.repo.ListSamples(runCtx, batchID, structToSampleFilter())
		if err != nil {
			return err
		}
		result, err := s.quality.Evaluate(runCtx, batch, protocol, samples)
		if err != nil {
			return err
		}
		_, err = s.alerts.OpenForQuality(runCtx, batch, result)
		return err
	}})
}

func (s *WorkerService) Running() bool { s.mu.Lock(); defer s.mu.Unlock(); return s.started }
