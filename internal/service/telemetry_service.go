package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
)

type TelemetryService struct {
	repo   *sqlite.Store
	mu     sync.RWMutex
	latest map[model.ID]model.AcousticSample
}

func NewTelemetryService(repo *sqlite.Store) *TelemetryService {
	return &TelemetryService{repo: repo, latest: make(map[model.ID]model.AcousticSample)}
}

func (s *TelemetryService) Record(ctx context.Context, sample model.AcousticSample) error {
	if err := sample.Validate(); err != nil {
		return err
	}
	if err := s.repo.AddSample(ctx, sample); err != nil {
		return fmt.Errorf("record telemetry: %w", err)
	}
	s.mu.Lock()
	s.latest[sample.SensorID] = sample
	s.mu.Unlock()
	return nil
}

func (s *TelemetryService) Latest(sensorID model.ID) (model.AcousticSample, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.latest[sensorID]
	return value, ok
}

func (s *TelemetryService) Stale(sensorID model.ID, now time.Time, threshold time.Duration) bool {
	value, ok := s.Latest(sensorID)
	return !ok || now.Sub(value.CapturedAt) > threshold
}
