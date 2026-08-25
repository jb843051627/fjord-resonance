package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/fjord-resonance/internal/ingest"
	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
)

type IngestService struct {
	samples  *SampleService
	pipeline *ingest.Pipeline
}

func NewIngestService(repo *sqlite.Store) *IngestService {
	service := &IngestService{}
	service.samples = NewSampleService(repo)
	service.pipeline = ingest.NewPipeline(service)
	return service
}

func (s *IngestService) Put(ctx context.Context, sample model.AcousticSample) error {
	return s.samples.Add(ctx, sample)
}

func (s *IngestService) Accept(ctx context.Context, frame ingest.Frame) error {
	if err := s.pipeline.Accept(ctx, frame); err != nil {
		return fmt.Errorf("accept acoustic frame: %w", err)
	}
	return nil
}

func (s *IngestService) AcceptMany(ctx context.Context, frames []ingest.Frame) error {
	for _, frame := range frames {
		if err := s.Accept(ctx, frame); err != nil {
			return err
		}
	}
	return nil
}

func (s *IngestService) Reset() { s.pipeline.Reset() }
