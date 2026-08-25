package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

func (s *Store) CreateProtocol(ctx context.Context, p model.Protocol) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if err := p.Validate(); err != nil {
		return fmt.Errorf("validate protocol: %w", err)
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO protocols(id,name,version,min_frequency,max_frequency,min_duration,max_duration,window_minutes,state,created_at) VALUES(?,?,?,?,?,?,?,?,?,?)`,
		p.ID, p.Name, p.Version, p.MinFrequencyHz, p.MaxFrequencyHz, p.MinDurationMS, p.MaxDurationMS, p.WindowMinutes, p.State, asText(p.CreatedAt))
	if err != nil {
		return fmt.Errorf("create protocol: %w", err)
	}
	return nil
}

func (s *Store) GetProtocol(ctx context.Context, id model.ID) (model.Protocol, error) {
	if err := checkContext(ctx); err != nil {
		return model.Protocol{}, err
	}
	var p model.Protocol
	var created string
	err := s.db.QueryRowContext(ctx, `SELECT id,name,version,min_frequency,max_frequency,min_duration,max_duration,window_minutes,state,created_at FROM protocols WHERE id=?`, id).
		Scan(&p.ID, &p.Name, &p.Version, &p.MinFrequencyHz, &p.MaxFrequencyHz, &p.MinDurationMS, &p.MaxDurationMS, &p.WindowMinutes, &p.State, &created)
	if err == sql.ErrNoRows {
		return model.Protocol{}, store.NotFound("protocol", string(id))
	}
	if err != nil {
		return model.Protocol{}, fmt.Errorf("get protocol: %w", err)
	}
	p.CreatedAt = fromText(created)
	return p, nil
}

func (s *Store) ListProtocols(ctx context.Context, state model.ProtocolState) ([]model.Protocol, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	query := `SELECT id,name,version,min_frequency,max_frequency,min_duration,max_duration,window_minutes,state,created_at FROM protocols`
	args := []any{}
	if state != "" {
		query += ` WHERE state=?`
		args = append(args, state)
	}
	query += ` ORDER BY name,version`
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list protocols: %w", err)
	}
	defer rows.Close()
	result := make([]model.Protocol, 0)
	for rows.Next() {
		var p model.Protocol
		var created string
		if err := rows.Scan(&p.ID, &p.Name, &p.Version, &p.MinFrequencyHz, &p.MaxFrequencyHz, &p.MinDurationMS, &p.MaxDurationMS, &p.WindowMinutes, &p.State, &created); err != nil {
			return nil, err
		}
		p.CreatedAt = fromText(created)
		result = append(result, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list protocols rows: %w", err)
	}
	return result, nil
}
