package model

import "strings"

type AlertRule struct {
	Code            string
	MinimumSeverity Severity
	Prefix          string
	Enabled         bool
}

func DefaultAlertRules() []AlertRule {
	return []AlertRule{{Code: "QUALITY_REJECTED", MinimumSeverity: SeverityCritical, Prefix: "quality", Enabled: true}, {Code: "QUALITY_REVIEW", MinimumSeverity: SeverityWarning, Prefix: "quality", Enabled: true}, {Code: "SENSOR_SILENT", MinimumSeverity: SeverityCritical, Prefix: "sensor", Enabled: true}}
}

func MatchRule(alert Alert, rules []AlertRule) (AlertRule, bool) {
	for _, rule := range rules {
		if rule.Enabled && rule.Code == alert.Code && (rule.Prefix == "" || strings.HasPrefix(strings.ToLower(alert.Message), strings.ToLower(rule.Prefix))) {
			return rule, true
		}
	}
	return AlertRule{}, false
}

func SeverityRank(value Severity) int {
	switch value {
	case SeverityInfo:
		return 1
	case SeverityWarning:
		return 2
	case SeverityCritical:
		return 3
	default:
		return 0
	}
}

func HigherSeverity(a, b Severity) Severity {
	if SeverityRank(a) >= SeverityRank(b) {
		return a
	}
	return b
}
