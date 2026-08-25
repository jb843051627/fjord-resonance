package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
)

type HealthService struct {
	repo      *sqlite.Store
	mu        sync.RWMutex
	states    map[model.ID]model.BuoyStatus
	refreshed time.Time
}

func NewHealthService(repo *sqlite.Store) *HealthService {
	return &HealthService{repo: repo, states: make(map[model.ID]model.BuoyStatus)}
}

func (s *HealthService) Refresh(ctx context.Context) error {
	buoys, err := s.repo.ListBuoys(ctx, structToBuoyFilter())
	if err != nil {
		return fmt.Errorf("refresh health: %w", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, buoy := range buoys {
		s.states[buoy.ID] = buoy.Status
	}
	s.refreshed = time.Now().UTC()
	return nil
}

func (s *HealthService) Status(id model.ID) (model.BuoyStatus, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.states[id]
	return value, ok
}

func (s *HealthService) Snapshot() map[model.ID]model.BuoyStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[model.ID]model.BuoyStatus, len(s.states))
	for key, value := range s.states {
		result[key] = value
	}
	return result
}

func (s *HealthService) LastRefresh() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.refreshed
}
