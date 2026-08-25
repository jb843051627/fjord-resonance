package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/fjord-resonance/internal/sqlite"
)

type MetricsService struct{ repo *sqlite.Store }

func NewMetricsService(repo *sqlite.Store) *MetricsService { return &MetricsService{repo: repo} }

func (s *MetricsService) Snapshot(ctx context.Context) (sqlite.HealthSnapshot, error) {
	snapshot, err := s.repo.HealthSnapshot(ctx)
	if err != nil {
		return sqlite.HealthSnapshot{}, fmt.Errorf("service metrics: %w", err)
	}
	return snapshot, nil
}

func (s *MetricsService) Ready(ctx context.Context) error {
	if !s.repo.DatabaseReady(ctx) {
		return fmt.Errorf("database is not ready")
	}
	return nil
}

func (s *MetricsService) Counts(ctx context.Context) (int, int, error) {
	open, err := s.repo.CountOpenAlerts(ctx)
	if err != nil {
		return 0, 0, err
	}
	active, err := s.repo.CountBuoysByStatus(ctx, "active")
	return active, open, err
}
