package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

type SensorService struct {
	repo   *sqlite.Store
	health map[model.ID]time.Time
	mu     sync.RWMutex
}

func NewSensorService(repo *sqlite.Store) *SensorService {
	return &SensorService{repo: repo, health: make(map[model.ID]time.Time)}
}

func (s *SensorService) Create(ctx context.Context, sensor model.Sensor) (model.Sensor, error) {
	if sensor.Status == "" {
		sensor.Status = model.SensorReady
	}
	if sensor.CreatedAt.IsZero() {
		sensor.CreatedAt = time.Now().UTC()
	}
	if err := sensor.Validate(); err != nil {
		return model.Sensor{}, fmt.Errorf("create sensor: %w", err)
	}
	if _, err := s.repo.GetBuoy(ctx, sensor.BuoyID); err != nil {
		return model.Sensor{}, fmt.Errorf("sensor buoy: %w", err)
	}
	if err := s.repo.CreateSensor(ctx, sensor); err != nil {
		return model.Sensor{}, err
	}
	return sensor, nil
}

func (s *SensorService) Get(ctx context.Context, id model.ID) (model.Sensor, error) {
	return s.repo.GetSensor(ctx, id)
}

func (s *SensorService) List(ctx context.Context, buoyID model.ID) ([]model.Sensor, error) {
	return s.repo.ListSensors(ctx, buoyID)
}

func (s *SensorService) Calibrate(ctx context.Context, id model.ID, factor float64) error {
	if factor <= 0 || factor > 100 {
		return store.ErrValidation
	}
	if _, err := s.Get(ctx, id); err != nil {
		return fmt.Errorf("calibrate sensor: %w", err)
	}
	return s.repo.UpdateSensorCalibration(ctx, id, factor, model.SensorReady)
}

func (s *SensorService) RecordHeartbeat(ctx context.Context, id model.ID, at time.Time) error {
	if _, err := s.Get(ctx, id); err != nil {
		return fmt.Errorf("record heartbeat: %w", err)
	}
	s.mu.RLock()
	s.health[id] = at
	s.mu.RUnlock()
	return nil
}

func (s *SensorService) LastHeartbeat(id model.ID) (time.Time, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.health[id]
	return value, ok
}

func (s *SensorService) MarkFault(ctx context.Context, id model.ID, reason string) error {
	if reason == "" {
		return store.ErrValidation
	}
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	return s.repo.UpdateSensorCalibration(ctx, id, 0, model.SensorFault)
}
