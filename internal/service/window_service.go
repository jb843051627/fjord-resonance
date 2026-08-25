package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
)

type WindowService struct{ repo *sqlite.Store }

func NewWindowService(repo *sqlite.Store) *WindowService { return &WindowService{repo: repo} }

func (s *WindowService) Overlaps(ctx context.Context, from, to time.Time) ([]model.CalibrationBatch, error) {
	if to.Before(from) {
		return nil, fmt.Errorf("invalid window")
	}
	batches, err := s.repo.BatchesInWindow(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("find overlapping batches: %w", err)
	}
	return batches, nil
}

func (s *WindowService) HasConflict(ctx context.Context, buoyID model.ID, from, to time.Time) (bool, error) {
	batches, err := s.Overlaps(ctx, from, to)
	if err != nil {
		return false, err
	}
	for _, batch := range batches {
		if batch.BuoyID == buoyID && !batch.Status.Terminal() {
			return true, nil
		}
	}
	return false, nil
}

func (s *WindowService) Validate(from, to time.Time) error {
	if from.IsZero() || to.IsZero() || !to.After(from) || to.Sub(from) > 24*time.Hour {
		return fmt.Errorf("window outside one-day operating limit")
	}
	return nil
}
