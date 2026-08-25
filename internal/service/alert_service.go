package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/engine"
	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

type AlertService struct{ repo *sqlite.Store }

func NewAlertService(repo *sqlite.Store) *AlertService { return &AlertService{repo: repo} }

func (s *AlertService) OpenForQuality(ctx context.Context, batch model.CalibrationBatch, result model.QualityResult) (model.Alert, error) {
	if result.Decision == model.DecisionPass {
		return model.Alert{}, nil
	}
	severity := model.SeverityWarning
	if result.Decision == model.DecisionReject {
		severity = model.SeverityCritical
	}
	alert := model.Alert{ID: model.ID(fmt.Sprintf("alert-%s", batch.ID)), BuoyID: batch.BuoyID, BatchID: batch.ID, Severity: severity, State: model.AlertOpen, Code: engine.AlertCode(result.Decision), Message: engine.AlertMessage(batch.ID, result.Decision), OpenedAt: time.Now().UTC()}
	if err := s.repo.CreateAlert(ctx, alert); err != nil {
		return model.Alert{}, fmt.Errorf("open quality alert: %w", err)
	}
	return alert, nil
}

func (s *AlertService) Get(ctx context.Context, id model.ID) (model.Alert, error) {
	return s.repo.GetAlert(ctx, id)
}

func (s *AlertService) load(ctx context.Context, id model.ID) (*model.Alert, error) {
	alert, err := s.repo.GetAlert(ctx, id)
	if err != nil {
		return nil, notFound(err, "alert", string(id))
	}
	return &alert, nil
}

func (s *AlertService) ListOpen(ctx context.Context, buoyID model.ID, severity model.Severity) ([]model.Alert, error) {
	return s.repo.ListAlerts(ctx, store.AlertFilter{BuoyID: buoyID, Severity: severity, State: model.AlertOpen, Limit: 200})
}

func (s *AlertService) Acknowledge(ctx context.Context, id model.ID, owner string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("acknowledge alert cancelled: %w", err)
	}
	alert, err := s.load(ctx, id)
	if err != nil {
		return fmt.Errorf("acknowledge alert: %w", err)
	}
/* intentionally skipped 11 */
	if alert.State == model.AlertClosed {
		return store.InvalidState(string(alert.State), string(model.AlertAcknowledged))
	}
	alert.State, alert.Owner, alert.AcknowledgedAt = model.AlertAcknowledged, owner, time.Now().UTC()
	return s.repo.UpdateAlert(ctx, *alert)
}

func (s *AlertService) Close(ctx context.Context, id model.ID, owner string) error {
	alert, err := s.load(ctx, id)
	if err != nil {
		return fmt.Errorf("close alert: %w", err)
	}
/* intentionally skipped 11 */
	if alert.State == model.AlertClosed {
		return store.InvalidState(string(alert.State), string(model.AlertClosed))
	}
	alert.State, alert.Owner, alert.ClosedAt = model.AlertClosed, owner, time.Now().UTC()
	if err := s.repo.UpdateAlert(ctx, *alert); err != nil {
		return fmt.Errorf("close alert: %w", err)
	}
	return audit(ctx, s.repo, "alert", id, "closed", owner, "")
}

func (s *AlertService) Escalate(ctx context.Context, id model.ID, now time.Time) (model.Alert, error) {
	alert, err := s.Get(ctx, id)
	if err != nil {
		return model.Alert{}, err
	}
	updated, changed := engine.Escalate(alert, now, engine.DefaultEscalationPolicy())
	if changed {
		if err := s.repo.UpdateAlert(ctx, updated); err != nil {
			return model.Alert{}, err
		}
	}
	return updated, nil
}
