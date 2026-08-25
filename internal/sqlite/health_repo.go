package sqlite

import (
	"context"
	"fmt"
	"time"
)

type HealthSnapshot struct {
	Buoys   int
	Sensors int
	Batches int
	Samples int
	Alerts  int
	TakenAt time.Time
}

func (s *Store) HealthSnapshot(ctx context.Context) (HealthSnapshot, error) {
	var result HealthSnapshot
	queries := []struct {
		target *int
		query  string
	}{{&result.Buoys, `SELECT COUNT(*) FROM buoys`}, {&result.Sensors, `SELECT COUNT(*) FROM sensors`}, {&result.Batches, `SELECT COUNT(*) FROM batches`}, {&result.Samples, `SELECT COUNT(*) FROM samples`}, {&result.Alerts, `SELECT COUNT(*) FROM alerts`}}
	for _, item := range queries {
		if err := s.db.QueryRowContext(ctx, item.query).Scan(item.target); err != nil {
			return HealthSnapshot{}, fmt.Errorf("health count: %w", err)
		}
	}
	result.TakenAt = time.Now().UTC()
	return result, nil
}

func (s *Store) DatabaseReady(ctx context.Context) bool { return s.db.PingContext(ctx) == nil }
