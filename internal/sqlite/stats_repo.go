package sqlite

import (
	"context"
	"fmt"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

type BatchStats struct {
	BatchID      model.ID
	Samples      int
	ValidSamples int
	Alerts       int
	OpenAlerts   int
}

func (s *Store) BatchStats(ctx context.Context, batchID model.ID) (BatchStats, error) {
	if err := checkContext(ctx); err != nil {
		return BatchStats{}, err
	}
	var stats BatchStats
	stats.BatchID = batchID
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*),COALESCE(SUM(valid),0) FROM samples WHERE batch_id=?`, batchID).Scan(&stats.Samples, &stats.ValidSamples); err != nil {
		return BatchStats{}, fmt.Errorf("batch sample stats: %w", err)
	}
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*),COALESCE(SUM(CASE WHEN state='open' THEN 1 ELSE 0 END),0) FROM alerts WHERE batch_id=?`, batchID).Scan(&stats.Alerts, &stats.OpenAlerts); err != nil {
		return BatchStats{}, fmt.Errorf("batch alert stats: %w", err)
	}
	return stats, nil
}

func (s *Store) CountBuoysByStatus(ctx context.Context, status model.BuoyStatus) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM buoys WHERE status=?`, status).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Store) CountOpenAlerts(ctx context.Context) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM alerts WHERE state='open'`).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
