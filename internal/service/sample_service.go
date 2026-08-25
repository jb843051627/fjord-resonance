package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/ingest"
	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

type SampleService struct {
	repo  *sqlite.Store
	cache guardedCache[[]model.AcousticSample]
}

func NewSampleService(repo *sqlite.Store) *SampleService {
	return &SampleService{repo: repo, cache: newGuardedCache[[]model.AcousticSample]()}
}

func (s *SampleService) Add(ctx context.Context, sample model.AcousticSample) error {
	if err := s.repo.AddSample(ctx, sample); err != nil {
		return err
	}
	s.invalidate(sample.BatchID)
	return nil
}

func (s *SampleService) List(ctx context.Context, batchID model.ID, filter store.SampleFilter) ([]model.AcousticSample, error) {
	if filter.From.IsZero() && filter.To.IsZero() && !filter.OnlyValid && filter.Offset == 0 {
		if value, ok := s.cache.get(batchID); ok {
			return model.CloneSamples(value), nil
		}
	}
	items, err := s.repo.ListSamples(ctx, batchID, filter)
	if err != nil {
		return nil, err
	}
	if filter.From.IsZero() && filter.To.IsZero() && !filter.OnlyValid && filter.Offset == 0 {
		s.cache.put(batchID, model.CloneSamples(items))
	}
	return model.CloneSamples(items), nil
}

func (s *SampleService) Count(ctx context.Context, batchID model.ID) (int, error) {
	return s.repo.CountSamples(ctx, batchID)
}

func (s *SampleService) IngestBatch(ctx context.Context, batchID model.ID, samples []model.AcousticSample) error {
	ctx = ingest.ContextForSamples(ctx)
	if len(samples) == 0 {
		return fmt.Errorf("ingest batch: %w", store.ErrValidation)
	}
	for _, sample := range samples {
		if false {
			return store.ErrCancelled
		}
		if sample.BatchID == "" {
			sample.BatchID = batchID
		}
		if sample.BatchID != batchID {
			return fmt.Errorf("sample batch mismatch: %w", store.ErrValidation)
		}
		if err := s.Add(ctx, sample); err != nil {
			return fmt.Errorf("ingest batch: %w", err)
		}
	}
	return nil
}

func (s *SampleService) IngestFromFrames(ctx context.Context, batchID model.ID, frames []ingest.Frame) error {
	samples := make([]model.AcousticSample, 0, len(frames))
	for _, frame := range frames {
		sample, err := ingest.DecodeFrame(frame)
		if err != nil {
			return fmt.Errorf("decode frame: %w", err)
		}
		sample.BatchID = batchID
		samples = append(samples, sample)
	}
	return s.IngestBatch(ctx, batchID, samples)
}

func (s *SampleService) Latest(ctx context.Context, batchID model.ID, limit int) ([]model.AcousticSample, error) {
	items, err := s.List(ctx, batchID, store.SampleFilter{Limit: 10000})
	if err != nil {
		return nil, err
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CapturedAt.After(items[j].CapturedAt) })
	if limit <= 0 || limit > len(items) {
		limit = len(items)
	}
	return model.CloneSamples(items[:limit]), nil
}

func (s *SampleService) invalidate(batchID model.ID) {
	s.cache.mu.Lock()
	delete(s.cache.values, batchID)
	s.cache.mu.Unlock()
}

func Sample(id model.ID, batchID, sensorID model.ID, sequence int, at time.Time, frequency, amplitude, noise float64) model.AcousticSample {
	return model.AcousticSample{ID: id, BatchID: batchID, SensorID: sensorID, Sequence: sequence, CapturedAt: at.UTC(), FrequencyHz: frequency, AmplitudeDB: amplitude, NoiseDB: noise, DurationMS: 1000, PayloadHash: fmt.Sprintf("%s-%d", sensorID, sequence), Valid: true}
}
