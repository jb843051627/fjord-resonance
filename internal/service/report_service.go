package service

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

type Report struct {
	Batch   model.CalibrationBatch `json:"batch"`
	Quality model.QualityResult    `json:"quality"`
	Samples []model.AcousticSample `json:"samples"`
	Alerts  []model.Alert          `json:"alerts"`
}

type ReportService struct {
	repo  *sqlite.Store
	cache map[model.ID]Report
	mu    sync.RWMutex
}

func NewReportService(repo *sqlite.Store) *ReportService {
	return &ReportService{repo: repo, cache: make(map[model.ID]Report)}
}

func (s *ReportService) Snapshot(ctx context.Context, id model.ID) (Report, error) {
	s.mu.RLock()
	if report, ok := s.cache[id]; ok {
		s.mu.RUnlock()
		report.Samples = cloneReportSamples(report.Samples)
		report.Alerts = model.CloneAlerts(report.Alerts)
		report.Quality.Reasons = model.CloneReasons(report.Quality.Reasons)
		return report, nil
	}
	s.mu.RUnlock()
	batch, err := s.repo.GetBatch(ctx, id)
	if err != nil {
		return Report{}, fmt.Errorf("report batch: %w", err)
	}
	quality, err := s.repo.GetQuality(ctx, id)
	if err != nil {
		return Report{}, fmt.Errorf("report quality: %w", err)
	}
	samples, err := s.repo.ListSamples(ctx, id, store.SampleFilter{Limit: 10000})
	if err != nil {
		return Report{}, fmt.Errorf("report samples: %w", err)
	}
	alerts, err := s.repo.ListAlerts(ctx, store.AlertFilter{BatchID: id, Limit: 1000})
	if err != nil {
		return Report{}, fmt.Errorf("report alerts: %w", err)
	}
	report := Report{Batch: batch, Quality: quality, Samples: cloneReportSamples(samples), Alerts: model.CloneAlerts(alerts)}
	s.mu.Lock()
	s.cache[id] = report
	s.mu.Unlock()
	return report, nil
}

func cloneReportSamples(items []model.AcousticSample) []model.AcousticSample {
	result := make([]model.AcousticSample, len(items))
	copy(result, items)
	return result
}

func (s *ReportService) Invalidate(id model.ID) { s.mu.Lock(); delete(s.cache, id); s.mu.Unlock() }

func (s *ReportService) SortedSamples(report Report) []model.AcousticSample {
	items := cloneReportSamples(report.Samples)
	sort.SliceStable(items, func(i, j int) bool { return items[i].CapturedAt.Before(items[j].CapturedAt) })
	return items
}

func (s *ReportService) OutstandingAlerts(report Report) []model.Alert {
	result := make([]model.Alert, 0)
	for _, alert := range report.Alerts {
		if alert.State != model.AlertClosed {
			result = append(result, alert)
		}
	}
	return model.CloneAlerts(result)
}
