package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

func (s *Store) CreateExport(ctx context.Context, job model.ExportJob) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if job.ID == "" || job.BatchID == "" || job.Format == "" {
		return store.ErrValidation
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO export_jobs(id,batch_id,format,state,requested_by,path,created_at,finished_at) VALUES(?,?,?,?,?,?,?,?)`,
		job.ID, job.BatchID, job.Format, job.State, job.RequestedBy, job.Path, asText(job.CreatedAt), asText(job.FinishedAt))
	if err != nil {
		return fmt.Errorf("create export: %w", err)
	}
	return nil
}

func (s *Store) GetExport(ctx context.Context, id model.ID) (model.ExportJob, error) {
	if err := checkContext(ctx); err != nil {
		return model.ExportJob{}, err
	}
	var job model.ExportJob
	var created, finished string
	err := s.db.QueryRowContext(ctx, `SELECT id,batch_id,format,state,requested_by,path,created_at,finished_at FROM export_jobs WHERE id=?`, id).
		Scan(&job.ID, &job.BatchID, &job.Format, &job.State, &job.RequestedBy, &job.Path, &created, &finished)
	if err == sql.ErrNoRows {
		return model.ExportJob{}, store.NotFound("export", string(id))
	}
	if err != nil {
		return model.ExportJob{}, fmt.Errorf("get export: %w", err)
	}
	job.CreatedAt, job.FinishedAt = fromText(created), fromText(finished)
	return job, nil
}

func (s *Store) UpdateExport(ctx context.Context, job model.ExportJob) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	result, err := s.db.ExecContext(ctx, `UPDATE export_jobs SET state=?,requested_by=?,path=?,finished_at=? WHERE id=?`, job.State, job.RequestedBy, job.Path, asText(job.FinishedAt), job.ID)
	if err != nil {
		return fmt.Errorf("update export: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return store.NotFound("export", string(job.ID))
	}
	return nil
}
