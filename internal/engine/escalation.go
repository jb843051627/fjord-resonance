package engine

import (
	"fmt"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

type EscalationPolicy struct {
	WarningAfter  time.Duration
	CriticalAfter time.Duration
}

func DefaultEscalationPolicy() EscalationPolicy {
	return EscalationPolicy{WarningAfter: 30 * time.Minute, CriticalAfter: 2 * time.Hour}
}

func Escalate(alert model.Alert, now time.Time, policy EscalationPolicy) (model.Alert, bool) {
	if alert.State == model.AlertClosed || alert.OpenedAt.IsZero() {
		return alert, false
	}
	age := now.Sub(alert.OpenedAt)
	changed := false
	if age >= policy.CriticalAfter && alert.Severity != model.SeverityCritical {
		alert.Severity, changed = model.SeverityCritical, true
	}
	if age >= policy.WarningAfter && alert.Severity == model.SeverityInfo {
		alert.Severity, changed = model.SeverityWarning, true
	}
	return alert, changed
}

func AlertCode(decision model.Decision) string {
	switch decision {
	case model.DecisionReject:
		return "QUALITY_REJECTED"
	case model.DecisionReview:
		return "QUALITY_REVIEW"
	default:
		return "QUALITY_ACCEPTED"
	}
}

func AlertMessage(batch model.ID, decision model.Decision) string {
	return fmt.Sprintf("batch %s requires %s handling", batch, decision)
}
