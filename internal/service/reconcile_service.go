package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

type ReconcileService struct{ repo *sqlite.Store }

func NewReconcileService(repo *sqlite.Store) *ReconcileService { return &ReconcileService{repo: repo} }

func (s *ReconcileService) Batch(ctx context.Context, id model.ID) (sqlite.BatchStats, error) {
	stats, err := s.repo.BatchStats(ctx, id)
	if err != nil {
		return sqlite.BatchStats{}, fmt.Errorf("reconcile batch: %w", err)
	}
	if stats.Samples > 0 && stats.ValidSamples == 0 {
		return stats, fmt.Errorf("batch has no valid samples: %w", store.ErrValidation)
	}
	return stats, nil
}

func (s *ReconcileService) RepairBuoy(ctx context.Context, id model.ID, seen time.Time) error {
	buoy, err := s.repo.GetBuoy(ctx, id)
	if err != nil {
		return err
	}
	if buoy.Status == model.BuoyRetired {
		return store.ErrInvalidState
	}
	return s.repo.UpdateBuoyStatus(ctx, id, model.BuoyActive, seen)
}

func (s *ReconcileService) RepairAlert(ctx context.Context, id model.ID, owner string) error {
	alert, err := s.repo.GetAlert(ctx, id)
	if err != nil {
		return err
	}
	if alert.State == model.AlertClosed {
		return store.ErrInvalidState
	}
	alert.Owner = owner
	return s.repo.UpdateAlert(ctx, alert)
}
