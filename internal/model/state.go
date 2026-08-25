package model

type BuoyStatus string
type SensorStatus string
type SensorKind string
type BatchStatus string
type ProtocolState string
type Decision string
type Severity string
type AlertState string
type ExportFormat string
type ExportState string

const (
	BuoyActive   BuoyStatus = "active"
	BuoyDrifting BuoyStatus = "drifting"
	BuoyOffline  BuoyStatus = "offline"
	BuoyRetired  BuoyStatus = "retired"

	SensorReady   SensorStatus = "ready"
	SensorWarmup  SensorStatus = "warmup"
	SensorFault   SensorStatus = "fault"
	SensorOffline SensorStatus = "offline"

	SensorHydrophone SensorKind = "hydrophone"
	SensorPressure   SensorKind = "pressure"
	SensorClock      SensorKind = "clock"

	BatchDraft     BatchStatus = "draft"
	BatchQueued    BatchStatus = "queued"
	BatchRunning   BatchStatus = "running"
	BatchEvaluated BatchStatus = "evaluated"
	BatchReview    BatchStatus = "review"
	BatchReleased  BatchStatus = "released"
	BatchRejected  BatchStatus = "rejected"
	BatchCancelled BatchStatus = "cancelled"

	ProtocolDraft   ProtocolState = "draft"
	ProtocolReady   ProtocolState = "ready"
	ProtocolRetired ProtocolState = "retired"

	DecisionPass   Decision = "pass"
	DecisionReview Decision = "review"
	DecisionReject Decision = "reject"

	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"

	AlertOpen         AlertState = "open"
	AlertAcknowledged AlertState = "acknowledged"
	AlertClosed       AlertState = "closed"

	ExportCSV   ExportFormat = "csv"
	ExportJSON  ExportFormat = "json"
	ExportReady ExportState  = "ready"
	ExportRun   ExportState  = "running"
	ExportDone  ExportState  = "done"
	ExportError ExportState  = "error"
)

func (s BatchStatus) Terminal() bool {
	return s == BatchReleased || s == BatchRejected || s == BatchCancelled
}

func (s AlertState) Closed() bool {
	return s == AlertClosed
}

func (d Decision) Accepted() bool {
	return d == DecisionPass || d == DecisionReview
}
