package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

func (s *Store) CreateSensor(ctx context.Context, sensor model.Sensor) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if err := sensor.Validate(); err != nil {
		return fmt.Errorf("validate sensor: %w", err)
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO sensors(id,buoy_id,serial,kind,status,sample_rate,calibration,last_reading,created_at) VALUES(?,?,?,?,?,?,?,?,?)`,
		sensor.ID, sensor.BuoyID, sensor.Serial, sensor.Kind, sensor.Status, sensor.SampleRate, sensor.Calibration, asText(sensor.LastReading), asText(sensor.CreatedAt))
	if err != nil {
		return fmt.Errorf("create sensor: %w", err)
	}
	return nil
}

func (s *Store) GetSensor(ctx context.Context, id model.ID) (model.Sensor, error) {
	if err := checkContext(ctx); err != nil {
		return model.Sensor{}, err
	}
	var sensor model.Sensor
	var lastReading, created string
	err := s.db.QueryRowContext(ctx, `SELECT id,buoy_id,serial,kind,status,sample_rate,calibration,last_reading,created_at FROM sensors WHERE id=?`, id).
		Scan(&sensor.ID, &sensor.BuoyID, &sensor.Serial, &sensor.Kind, &sensor.Status, &sensor.SampleRate, &sensor.Calibration, &lastReading, &created)
	if err == sql.ErrNoRows {
		return model.Sensor{}, store.NotFound("sensor", string(id))
	}
	if err != nil {
		return model.Sensor{}, fmt.Errorf("get sensor: %w", err)
	}
	sensor.LastReading, sensor.CreatedAt = fromText(lastReading), fromText(created)
	return sensor, nil
}

func (s *Store) ListSensors(ctx context.Context, buoyID model.ID) ([]model.Sensor, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id,buoy_id,serial,kind,status,sample_rate,calibration,last_reading,created_at FROM sensors WHERE buoy_id=? ORDER BY serial`, buoyID)
	if err != nil {
		return nil, fmt.Errorf("list sensors: %w", err)
	}
	defer rows.Close()
	result := make([]model.Sensor, 0)
	for rows.Next() {
		var sensor model.Sensor
		var lastReading, created string
		if err := rows.Scan(&sensor.ID, &sensor.BuoyID, &sensor.Serial, &sensor.Kind, &sensor.Status, &sensor.SampleRate, &sensor.Calibration, &lastReading, &created); err != nil {
			return nil, err
		}
		sensor.LastReading, sensor.CreatedAt = fromText(lastReading), fromText(created)
		result = append(result, sensor)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list sensors rows: %w", err)
	}
	return result, nil
}

func (s *Store) UpdateSensorCalibration(ctx context.Context, id model.ID, calibration float64, status model.SensorStatus) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if calibration < 0 || calibration > 100 {
		return store.ErrValidation
	}
	result, err := s.db.ExecContext(ctx, `UPDATE sensors SET calibration=?,status=?,last_reading=? WHERE id=?`, calibration, status, asText(time.Now()), id)
	if err != nil {
		return fmt.Errorf("update sensor calibration: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return store.NotFound("sensor", string(id))
	}
	return nil
}
