package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

type BuoyService struct {
	repo  *sqlite.Store
	cache guardedCache[model.Buoy]
}

func NewBuoyService(repo *sqlite.Store) *BuoyService {
	return &BuoyService{repo: repo, cache: newGuardedCache[model.Buoy]()}
}

func (s *BuoyService) Create(ctx context.Context, buoy model.Buoy) (model.Buoy, error) {
	if buoy.Status == "" {
		buoy.Status = model.BuoyActive
	}
	if buoy.CreatedAt.IsZero() {
		buoy.CreatedAt = time.Now().UTC()
	}
	buoy.UpdatedAt = buoy.CreatedAt
	if err := buoy.Validate(); err != nil {
		return model.Buoy{}, fmt.Errorf("create buoy: %w", err)
	}
	if err := s.repo.CreateBuoy(ctx, buoy); err != nil {
		return model.Buoy{}, err
	}
	s.cache.put(buoy.ID, buoy)
	_ = audit(ctx, s.repo, "buoy", buoy.ID, "created", "system", buoy.Name)
	return buoy, nil
}

func (s *BuoyService) Get(ctx context.Context, id model.ID) (model.Buoy, error) {
	if buoy, ok := s.cache.get(id); ok {
		return buoy, nil
	}
	buoy, err := s.repo.GetBuoy(ctx, id)
	if err != nil {
		return model.Buoy{}, notFound(err, "buoy", string(id))
	}
	s.cache.put(id, buoy)
	return buoy, nil
}

func (s *BuoyService) Require(ctx context.Context, id model.ID) (model.Buoy, error) {
	buoy, err := s.repo.GetBuoy(ctx, id)
	if err != nil {
		return model.Buoy{}, err
	}
	return buoy, nil
}

func (s *BuoyService) List(ctx context.Context, status model.BuoyStatus, name string, limit int) ([]model.Buoy, error) {
	return s.repo.ListBuoys(ctx, store.BuoyFilter{Status: status, Name: strings.TrimSpace(name), Limit: limit})
}

func (s *BuoyService) MarkSeen(ctx context.Context, id model.ID, at time.Time) error {
	buoy, err := s.Require(ctx, id)
	if err != nil {
		return fmt.Errorf("mark buoy seen: %w", err)
	}
	if buoy.Status == model.BuoyRetired {
		return fmt.Errorf("retired buoy: %w", store.ErrInvalidState)
	}
	if err := s.repo.UpdateBuoyStatus(ctx, id, model.BuoyActive, at); err != nil {
		return fmt.Errorf("mark buoy seen: %w", err)
	}
	buoy.Status, buoy.LastSeen, buoy.UpdatedAt = model.BuoyActive, at, time.Now().UTC()
	s.cache.put(id, buoy)
	return nil
}

func (s *BuoyService) SetStatus(ctx context.Context, id model.ID, status model.BuoyStatus) error {
	if _, err := s.Require(ctx, id); err != nil {
		return err
	}
	if status == "" {
		return store.ErrValidation
	}
	if err := s.repo.UpdateBuoyStatus(ctx, id, status, time.Now().UTC()); err != nil {
		return err
	}
	s.cache.mu.Lock()
	if buoy, ok := s.cache.values[id]; ok {
		buoy.Status = status
		s.cache.values[id] = buoy
	}
	s.cache.mu.Unlock()
	return nil
}
