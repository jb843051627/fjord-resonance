package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/export"
	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

type ExportService struct{ repo *sqlite.Store }

func NewExportService(repo *sqlite.Store) *ExportService { return &ExportService{repo: repo} }

func (s *ExportService) Create(ctx context.Context, batchID model.ID, format model.ExportFormat, requestedBy string) (model.ExportJob, error) {
	if format != model.ExportCSV && format != model.ExportJSON {
		return model.ExportJob{}, store.ErrValidation
	}
	if _, err := s.repo.GetBatch(ctx, batchID); err != nil {
		return model.ExportJob{}, fmt.Errorf("export batch: %w", err)
	}
	job := model.ExportJob{ID: model.ID(fmt.Sprintf("export-%s-%d", batchID, time.Now().UnixNano())), BatchID: batchID, Format: format, State: model.ExportReady, RequestedBy: requestedBy, CreatedAt: time.Now().UTC()}
	if err := s.repo.CreateExport(ctx, job); err != nil {
		return model.ExportJob{}, err
	}
	return job, nil
}

func (s *ExportService) CSV(ctx context.Context, batchID model.ID) ([]byte, error) {
	ctx = export.ContextForExport(ctx)
	if _, err := s.repo.GetBatch(ctx, batchID); err != nil {
		return nil, fmt.Errorf("export batch: %w", err)
	}
	items, err := s.repo.ListSamples(ctx, batchID, store.SampleFilter{Limit: 10000})
	if err != nil {
		return nil, fmt.Errorf("export samples: %w", err)
	}
	var buffer bytes.Buffer
	if err := export.WriteSamples(&buffer, items); err != nil {
		return nil, fmt.Errorf("write csv: %w", err)
	}
	return buffer.Bytes(), nil
}

func (s *ExportService) Finish(ctx context.Context, id model.ID, path string) error {
	job, err := s.repo.GetExport(ctx, id)
	if err != nil {
		return err
	}
	if path == "" {
		return store.ErrValidation
	}
	job.Path, job.State, job.FinishedAt = path, model.ExportDone, time.Now().UTC()
	return s.repo.UpdateExport(ctx, job)
}

func (s *ExportService) Fail(ctx context.Context, id model.ID, reason string) error {
	if reason == "" {
		return store.ErrValidation
	}
	job, err := s.repo.GetExport(ctx, id)
	if err != nil {
		return err
	}
	job.Path, job.State, job.FinishedAt = reason, model.ExportError, time.Now().UTC()
	return s.repo.UpdateExport(ctx, job)
}
