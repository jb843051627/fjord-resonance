package sqlite

import (
	"context"
	"fmt"
	"time"
)

func (s *Store) PruneAudit(ctx context.Context, before time.Time) (int64, error) {
	if before.IsZero() {
		return 0, fmt.Errorf("prune time is empty")
	}
	result, err := s.db.ExecContext(ctx, `DELETE FROM audit_events WHERE created_at < ?`, asText(before))
	if err != nil {
		return 0, fmt.Errorf("prune audit: %w", err)
	}
	return result.RowsAffected()
}

func (s *Store) Vacuum(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `VACUUM`)
	return err
}

func (s *Store) Ping(ctx context.Context) error { return s.db.PingContext(ctx) }
