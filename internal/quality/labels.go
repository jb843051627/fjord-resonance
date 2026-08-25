package quality

import "github.com/jb843051627/fjord-resonance/internal/model"

func DecisionLabel(decision model.Decision) string {
	switch decision {
	case model.DecisionPass:
		return "accepted"
	case model.DecisionReview:
		return "manual review"
	case model.DecisionReject:
		return "rejected"
	default:
		return "unknown"
	}
}

func SeverityForDecision(decision model.Decision) model.Severity {
	if decision == model.DecisionReject {
		return model.SeverityCritical
	}
	if decision == model.DecisionReview {
		return model.SeverityWarning
	}
	return model.SeverityInfo
}

func ShouldOpenAlert(decision model.Decision) bool { return decision != model.DecisionPass }
