package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

type MaintenanceService struct{ repo *sqlite.Store }

func NewMaintenanceService(repo *sqlite.Store) *MaintenanceService {
	return &MaintenanceService{repo: repo}
}

func (s *MaintenanceService) ExpireBatches(ctx context.Context, now time.Time) ([]model.ID, error) {
	ids, err := s.repo.ExpiredBatches(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("find expired batches: %w", err)
	}
	for _, id := range ids {
		batch, err := s.repo.GetBatch(ctx, id)
		if err != nil {
			return nil, err
		}
		batch.Status, batch.Summary = model.BatchCancelled, "window expired"
		if err := s.repo.UpdateBatch(ctx, batch); err != nil {
			return nil, err
		}
	}
	return ids, nil
}

func (s *MaintenanceService) MarkOffline(ctx context.Context, id model.ID, at time.Time) error {
	if err := s.repo.UpdateBuoyStatus(ctx, id, model.BuoyOffline, at); err != nil {
		return fmt.Errorf("mark offline: %w", err)
	}
	return nil
}

func (s *MaintenanceService) ValidateRepair(status model.SensorStatus) error {
	if status != model.SensorFault {
		return store.InvalidState(string(status), string(model.SensorReady))
	}
	return nil
}

func (s *MaintenanceService) ReopenAlert(ctx context.Context, id model.ID) error {
	alert, err := s.repo.GetAlert(ctx, id)
	if err != nil {
		return err
	}
	if alert.State != model.AlertClosed {
		return store.ErrInvalidState
	}
	alert.State, alert.ClosedAt = model.AlertOpen, time.Time{}
	return s.repo.UpdateAlert(ctx, alert)
}
