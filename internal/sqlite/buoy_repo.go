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

func (s *Store) CreateBuoy(ctx context.Context, b model.Buoy) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if err := b.Validate(); err != nil {
		return fmt.Errorf("validate buoy: %w", err)
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO buoys(id,name,latitude,longitude,depth_meters,status,last_seen,notes,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?)`,
		b.ID, b.Name, b.Latitude, b.Longitude, b.DepthMeters, b.Status, asText(b.LastSeen), b.Notes, asText(b.CreatedAt), asText(b.UpdatedAt))
	if err != nil {
		return fmt.Errorf("create buoy: %w", err)
	}
	return nil
}

func (s *Store) GetBuoy(ctx context.Context, id model.ID) (model.Buoy, error) {
	if err := checkContext(ctx); err != nil {
		return model.Buoy{}, err
	}
	var b model.Buoy
	var seenText, created, updated string
	err := s.db.QueryRowContext(ctx, `SELECT id,name,latitude,longitude,depth_meters,status,last_seen,notes,created_at,updated_at FROM buoys WHERE id=?`, id).
		Scan(&b.ID, &b.Name, &b.Latitude, &b.Longitude, &b.DepthMeters, &b.Status, &seenText, &b.Notes, &created, &updated)
	if err == sql.ErrNoRows {
		return model.Buoy{}, notFound("buoy", string(id))
	}
	if err != nil {
		return model.Buoy{}, fmt.Errorf("get buoy: %w", err)
	}
	b.LastSeen, b.CreatedAt, b.UpdatedAt = fromText(seenText), fromText(created), fromText(updated)
	return b, nil
}

func (s *Store) ListBuoys(ctx context.Context, filter store.BuoyFilter) ([]model.Buoy, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	query := `SELECT id,name,latitude,longitude,depth_meters,status,last_seen,notes,created_at,updated_at FROM buoys WHERE 1=1`
	args := make([]any, 0, 4)
	if filter.Status != "" {
		query += ` AND status=?`
		args = append(args, filter.Status)
	}
	if strings.TrimSpace(filter.Name) != "" {
		query += ` AND name LIKE ?`
		args = append(args, "%"+strings.TrimSpace(filter.Name)+"%")
	}
	query += ` ORDER BY name LIMIT ? OFFSET ?`
	args = append(args, store.NormalizeLimit(filter.Limit, 50, 200), max(filter.Offset, 0))
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list buoys: %w", err)
	}
	defer rows.Close()
	result := make([]model.Buoy, 0)
	for rows.Next() {
		var b model.Buoy
		var seenText, created, updated string
		if err := rows.Scan(&b.ID, &b.Name, &b.Latitude, &b.Longitude, &b.DepthMeters, &b.Status, &seenText, &b.Notes, &created, &updated); err != nil {
			return nil, err
		}
		b.LastSeen, b.CreatedAt, b.UpdatedAt = fromText(seenText), fromText(created), fromText(updated)
		result = append(result, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list buoys rows: %w", err)
	}
	return result, nil
}

func (s *Store) UpdateBuoyStatus(ctx context.Context, id model.ID, status model.BuoyStatus, seen time.Time) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	result, err := s.db.ExecContext(ctx, `UPDATE buoys SET status=?,last_seen=?,updated_at=? WHERE id=?`, status, asText(seen), asText(time.Now()), id)
	if err != nil {
		return fmt.Errorf("update buoy status: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return store.NotFound("buoy", string(id))
	}
	return nil
}
