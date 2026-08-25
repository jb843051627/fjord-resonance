package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

type AlertQueryService struct{ repo *sqlite.Store }

func NewAlertQueryService(repo *sqlite.Store) *AlertQueryService {
	return &AlertQueryService{repo: repo}
}

func (s *AlertQueryService) ForBuoy(ctx context.Context, buoyID model.ID, limit int) ([]model.Alert, error) {
	if buoyID == "" {
		return nil, store.ErrValidation
	}
	items, err := s.repo.ListAlerts(ctx, store.AlertFilter{BuoyID: buoyID, Limit: limit})
	if err != nil {
		return nil, fmt.Errorf("alerts for buoy: %w", err)
	}
	return model.CloneAlerts(items), nil
}

func (s *AlertQueryService) CriticalFirst(ctx context.Context, buoyID model.ID) ([]model.Alert, error) {
	items, err := s.ForBuoy(ctx, buoyID, 1000)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(items, func(i, j int) bool {
		return model.SeverityRank(items[i].Severity) > model.SeverityRank(items[j].Severity)
	})
	return items, nil
}

func (s *AlertQueryService) OpenCount(ctx context.Context) (int, error) {
	return s.repo.CountOpenAlerts(ctx)
}
