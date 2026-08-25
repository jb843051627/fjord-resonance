package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

func (s *Store) CreateBatch(ctx context.Context, b model.CalibrationBatch) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if b.ID == "" || b.BuoyID == "" || b.ProtocolID == "" || !b.WindowEnd.After(b.WindowStart) {
		return store.ErrValidation
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO batches(id,buoy_id,protocol_id,status,window_start,window_end,started_at,completed_at,reviewer,summary,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`,
		b.ID, b.BuoyID, b.ProtocolID, b.Status, asText(b.WindowStart), asText(b.WindowEnd), asText(b.StartedAt), asText(b.CompletedAt), b.Reviewer, b.Summary, asText(b.CreatedAt), asText(b.UpdatedAt))
	if err != nil {
		return fmt.Errorf("create batch: %w", err)
	}
	return nil
}

func (s *Store) GetBatch(ctx context.Context, id model.ID) (model.CalibrationBatch, error) {
	if err := checkContext(ctx); err != nil {
		return model.CalibrationBatch{}, err
	}
	var b model.CalibrationBatch
	var start, end, started, completed, created, updated string
	err := s.db.QueryRowContext(ctx, `SELECT id,buoy_id,protocol_id,status,window_start,window_end,started_at,completed_at,reviewer,summary,created_at,updated_at FROM batches WHERE id=?`, id).
		Scan(&b.ID, &b.BuoyID, &b.ProtocolID, &b.Status, &start, &end, &started, &completed, &b.Reviewer, &b.Summary, &created, &updated)
	if err == sql.ErrNoRows {
		return model.CalibrationBatch{}, store.NotFound("batch", string(id))
	}
	if err != nil {
		return model.CalibrationBatch{}, fmt.Errorf("get batch: %w", err)
	}
	b.WindowStart, b.WindowEnd = fromText(start), fromText(end)
	b.StartedAt, b.CompletedAt = fromText(started), fromText(completed)
	b.CreatedAt, b.UpdatedAt = fromText(created), fromText(updated)
	return b, nil
}

func (s *Store) ListBatches(ctx context.Context, filter store.BatchFilter) ([]model.CalibrationBatch, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	query := `SELECT id,buoy_id,protocol_id,status,window_start,window_end,started_at,completed_at,reviewer,summary,created_at,updated_at FROM batches WHERE 1=1`
	args := make([]any, 0, 8)
	if filter.BuoyID != "" {
		query += ` AND buoy_id=?`
		args = append(args, filter.BuoyID)
	}
	if filter.Status != "" {
		query += ` AND status=?`
		args = append(args, filter.Status)
	}
	if !filter.From.IsZero() {
		query += ` AND window_start>=?`
		args = append(args, asText(filter.From))
	}
	if !filter.To.IsZero() {
		query += ` AND window_end<=?`
		args = append(args, asText(filter.To))
	}
	query += ` ORDER BY window_start DESC LIMIT ? OFFSET ?`
	args = append(args, store.NormalizeLimit(filter.Limit, 50, 200), max(filter.Offset, 0))
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list batches: %w", err)
	}
	defer rows.Close()
	result := make([]model.CalibrationBatch, 0)
	for rows.Next() {
		var b model.CalibrationBatch
		var start, end, started, completed, created, updated string
		if err := rows.Scan(&b.ID, &b.BuoyID, &b.ProtocolID, &b.Status, &start, &end, &started, &completed, &b.Reviewer, &b.Summary, &created, &updated); err != nil {
			return nil, err
		}
		b.WindowStart, b.WindowEnd = fromText(start), fromText(end)
		b.StartedAt, b.CompletedAt = fromText(started), fromText(completed)
		b.CreatedAt, b.UpdatedAt = fromText(created), fromText(updated)
		result = append(result, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list batches rows: %w", err)
	}
	return result, nil
}

func (s *Store) UpdateBatch(ctx context.Context, b model.CalibrationBatch) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	result, err := s.db.ExecContext(ctx, `UPDATE batches SET status=?,window_start=?,window_end=?,started_at=?,completed_at=?,reviewer=?,summary=?,updated_at=? WHERE id=?`,
		b.Status, asText(b.WindowStart), asText(b.WindowEnd), asText(b.StartedAt), asText(b.CompletedAt), b.Reviewer, strings.TrimSpace(b.Summary), asText(time.Now()), b.ID)
	if err != nil {
		return fmt.Errorf("update batch: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return store.NotFound("batch", string(b.ID))
	}
	return nil
}
