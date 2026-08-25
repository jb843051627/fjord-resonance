package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

func (s *Store) BatchesInWindow(ctx context.Context, from, to time.Time) ([]model.CalibrationBatch, error) {
	if to.Before(from) {
		return nil, fmt.Errorf("window end before start")
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id,buoy_id,protocol_id,status,window_start,window_end,started_at,completed_at,reviewer,summary,created_at,updated_at FROM batches WHERE window_start < ? AND window_end > ? ORDER BY window_start`, asText(to), asText(from))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]model.CalibrationBatch, 0)
	for rows.Next() {
		var batch model.CalibrationBatch
		var start, end, started, completed, created, updated string
		if err := rows.Scan(&batch.ID, &batch.BuoyID, &batch.ProtocolID, &batch.Status, &start, &end, &started, &completed, &batch.Reviewer, &batch.Summary, &created, &updated); err != nil {
			return nil, err
		}
		batch.WindowStart, batch.WindowEnd, batch.StartedAt, batch.CompletedAt = fromText(start), fromText(end), fromText(started), fromText(completed)
		batch.CreatedAt, batch.UpdatedAt = fromText(created), fromText(updated)
		result = append(result, batch)
	}
	return result, rows.Err()
}

func (s *Store) ExpiredBatches(ctx context.Context, now time.Time) ([]model.ID, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id FROM batches WHERE status IN ('queued','running') AND window_end < ?`, asText(now))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]model.ID, 0)
	for rows.Next() {
		var id model.ID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		result = append(result, id)
	}
	return result, rows.Err()
}
