package sqlite

import (
	"context"
	"fmt"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

func (s *Store) AppendAudit(ctx context.Context, event model.AuditEvent) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if event.ID == "" || event.Entity == "" || event.EntityID == "" || event.Action == "" {
		return fmt.Errorf("audit event is incomplete")
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO audit_events(id,entity,entity_id,action,actor,details,created_at) VALUES(?,?,?,?,?,?,?)`,
		event.ID, event.Entity, event.EntityID, event.Action, event.Actor, event.Details, asText(event.CreatedAt))
	if err != nil {
		return fmt.Errorf("append audit: %w", err)
	}
	return nil
}

func (s *Store) ListAudit(ctx context.Context, entity string, entityID model.ID) ([]model.AuditEvent, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id,entity,entity_id,action,actor,details,created_at FROM audit_events WHERE entity=? AND entity_id=? ORDER BY created_at`, entity, entityID)
	if err != nil {
		return nil, fmt.Errorf("list audit: %w", err)
	}
	defer rows.Close()
	result := make([]model.AuditEvent, 0)
	for rows.Next() {
		var event model.AuditEvent
		var created string
		if err := rows.Scan(&event.ID, &event.Entity, &event.EntityID, &event.Action, &event.Actor, &event.Details, &created); err != nil {
			return nil, err
		}
		event.CreatedAt = fromText(created)
		result = append(result, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list audit rows: %w", err)
	}
	return result, nil
}
