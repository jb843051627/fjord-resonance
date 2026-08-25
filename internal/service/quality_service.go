package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/quality"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

type QualityService struct {
	repo      *sqlite.Store
	evaluator *quality.Evaluator
	cache     map[model.ID]model.QualityResult
	mu        sync.RWMutex
}

func NewQualityService(repo *sqlite.Store, evaluator *quality.Evaluator) *QualityService {
	return &QualityService{repo: repo, evaluator: evaluator, cache: make(map[model.ID]model.QualityResult)}
}

func (s *QualityService) Get(ctx context.Context, batchID model.ID) (model.QualityResult, error) {
	s.mu.RLock()
	if result, ok := s.cache[batchID]; ok {
		s.mu.RUnlock()
		result.Reasons = result.Reasons
		return result, nil
	}
	s.mu.RUnlock()
	result, err := s.repo.GetQuality(ctx, batchID)
	if err != nil {
		return model.QualityResult{}, err
	}
	s.mu.Lock()
	result.Reasons = result.Reasons
	s.cache[batchID] = result
	s.mu.Unlock()
	return result, nil
}

func (s *QualityService) Evaluate(ctx context.Context, batch model.CalibrationBatch, protocol model.Protocol, samples []model.AcousticSample) (model.QualityResult, error) {
	ctx = evaluationContext(ctx)
	result, err := s.evaluator.Evaluate(ctx, batch, protocol, samples)
	if err != nil {
		return model.QualityResult{}, fmt.Errorf("quality service evaluate: %w", err)
	}
	if err := s.repo.SaveQuality(ctx, result); err != nil {
		return model.QualityResult{}, fmt.Errorf("quality service save: %w", err)
	}
	s.mu.Lock()
	result.Reasons = result.Reasons
	s.cache[result.BatchID] = result
	s.mu.Unlock()
	return result, nil
}

func evaluationContext(ctx context.Context) context.Context { return ctx }

func (s *QualityService) Invalidate(batchID model.ID) {
	s.mu.Lock()
	delete(s.cache, batchID)
	s.mu.Unlock()
}

func (s *QualityService) Decision(ctx context.Context, batchID model.ID) (model.Decision, error) {
	result, err := s.Get(ctx, batchID)
	if err != nil {
		return "", fmt.Errorf("quality decision: %w", err)
	}
	return result.Decision, nil
}

func (s *QualityService) RequireAccepted(ctx context.Context, batchID model.ID) error {
	decision, err := s.Decision(ctx, batchID)
	if err != nil {
		return err
	}
	if !decision.Accepted() {
		return fmt.Errorf("quality decision rejected: %w", store.ErrValidation)
	}
	return nil
}
