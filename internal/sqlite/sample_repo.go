package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

func (s *Store) AddSample(ctx context.Context, sample model.AcousticSample) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if err := sample.Validate(); err != nil {
		return fmt.Errorf("validate sample: %w", err)
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO samples(id,batch_id,sensor_id,captured_at,frequency,amplitude,noise,duration_ms,sequence_no,payload_hash,valid) VALUES(?,?,?,?,?,?,?,?,?,?,?)`,
		sample.ID, sample.BatchID, sample.SensorID, asText(sample.CapturedAt), sample.FrequencyHz, sample.AmplitudeDB, sample.NoiseDB, sample.DurationMS, sample.Sequence, sample.PayloadHash, boolInt(sample.Valid))
	if err != nil {
		return fmt.Errorf("add sample: %w", err)
	}
	return nil
}

func (s *Store) ListSamples(ctx context.Context, batchID model.ID, filter store.SampleFilter) ([]model.AcousticSample, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	query := `SELECT id,batch_id,sensor_id,captured_at,frequency,amplitude,noise,duration_ms,sequence_no,payload_hash,valid FROM samples WHERE batch_id=?`
	args := []any{batchID}
	if !filter.From.IsZero() {
		query += ` AND captured_at>=?`
		args = append(args, asText(filter.From))
	}
	if !filter.To.IsZero() {
		query += ` AND captured_at<=?`
		args = append(args, asText(filter.To))
	}
	if filter.OnlyValid {
		query += ` AND valid=1`
	}
	query += ` ORDER BY sequence_no,captured_at LIMIT ? OFFSET ?`
	args = append(args, store.NormalizeLimit(filter.Limit, 1000, 10000), max(filter.Offset, 0))
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list samples: %w", err)
	}
	defer rows.Close()
	result := make([]model.AcousticSample, 0)
	for rows.Next() {
		var sample model.AcousticSample
		var captured string
		var valid int
		if err := rows.Scan(&sample.ID, &sample.BatchID, &sample.SensorID, &captured, &sample.FrequencyHz, &sample.AmplitudeDB, &sample.NoiseDB, &sample.DurationMS, &sample.Sequence, &sample.PayloadHash, &valid); err != nil {
			return nil, err
		}
		sample.CapturedAt, sample.Valid = fromText(captured), valid != 0
		result = append(result, sample)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list samples rows: %w", err)
	}
	return result, nil
}

func (s *Store) CountSamples(ctx context.Context, batchID model.ID) (int, error) {
	if err := checkContext(ctx); err != nil {
		return 0, err
	}
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM samples WHERE batch_id=?`, batchID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count samples: %w", err)
	}
	return count, nil
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

var _ = sql.ErrNoRows
