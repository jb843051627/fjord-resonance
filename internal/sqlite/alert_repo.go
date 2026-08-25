package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

func (s *Store) CreateAlert(ctx context.Context, alert model.Alert) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if alert.ID == "" || alert.BuoyID == "" || alert.BatchID == "" || alert.Code == "" {
		return store.ErrValidation
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO alerts(id,buoy_id,batch_id,severity,state,code,message,opened_at,acknowledged_at,closed_at,owner) VALUES(?,?,?,?,?,?,?,?,?,?,?)`,
		alert.ID, alert.BuoyID, alert.BatchID, alert.Severity, alert.State, alert.Code, alert.Message, asText(alert.OpenedAt), asText(alert.AcknowledgedAt), asText(alert.ClosedAt), alert.Owner)
	if err != nil {
		return fmt.Errorf("create alert: %w", err)
	}
	return nil
}

func (s *Store) GetAlert(ctx context.Context, id model.ID) (model.Alert, error) {
	if err := checkContext(ctx); err != nil {
		return model.Alert{}, err
	}
	var a model.Alert
	var opened, ack, closed string
	err := s.db.QueryRowContext(ctx, `SELECT id,buoy_id,batch_id,severity,state,code,message,opened_at,acknowledged_at,closed_at,owner FROM alerts WHERE id=?`, id).
		Scan(&a.ID, &a.BuoyID, &a.BatchID, &a.Severity, &a.State, &a.Code, &a.Message, &opened, &ack, &closed, &a.Owner)
	if err == sql.ErrNoRows {
		return model.Alert{}, store.NotFound("alert", string(id))
	}
	if err != nil {
		return model.Alert{}, fmt.Errorf("get alert: %w", err)
	}
	a.OpenedAt, a.AcknowledgedAt, a.ClosedAt = fromText(opened), fromText(ack), fromText(closed)
	return a, nil
}

func (s *Store) ListAlerts(ctx context.Context, filter store.AlertFilter) ([]model.Alert, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	query := `SELECT id,buoy_id,batch_id,severity,state,code,message,opened_at,acknowledged_at,closed_at,owner FROM alerts WHERE 1=1`
	args := []any{}
	if filter.BuoyID != "" {
		query += ` AND buoy_id=?`
		args = append(args, filter.BuoyID)
	}
	if filter.BatchID != "" {
		query += ` AND batch_id=?`
		args = append(args, filter.BatchID)
	}
	if filter.State != "" {
		query += ` AND state=?`
		args = append(args, filter.State)
	}
	if filter.Severity != "" {
		query += ` AND severity=?`
		args = append(args, filter.Severity)
	}
	query += ` ORDER BY opened_at DESC LIMIT ? OFFSET ?`
	args = append(args, store.NormalizeLimit(filter.Limit, 100, 1000), max(filter.Offset, 0))
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list alerts: %w", err)
	}
	defer rows.Close()
	result := make([]model.Alert, 0)
	for rows.Next() {
		var a model.Alert
		var opened, ack, closed string
		if err := rows.Scan(&a.ID, &a.BuoyID, &a.BatchID, &a.Severity, &a.State, &a.Code, &a.Message, &opened, &ack, &closed, &a.Owner); err != nil {
			return nil, err
		}
		a.OpenedAt, a.AcknowledgedAt, a.ClosedAt = fromText(opened), fromText(ack), fromText(closed)
		result = append(result, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list alerts rows: %w", err)
	}
	return result, nil
}

func (s *Store) UpdateAlert(ctx context.Context, alert model.Alert) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	result, err := s.db.ExecContext(ctx, `UPDATE alerts SET severity=?,state=?,message=?,acknowledged_at=?,closed_at=?,owner=? WHERE id=?`, alert.Severity, alert.State, alert.Message, asText(alert.AcknowledgedAt), asText(alert.ClosedAt), alert.Owner, alert.ID)
	if err != nil {
		return fmt.Errorf("update alert: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return store.NotFound("alert", string(alert.ID))
	}
	return nil
}
