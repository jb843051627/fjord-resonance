package model

import "time"

type ID string

type Buoy struct {
	ID          ID         `json:"id"`
	Name        string     `json:"name"`
	Latitude    float64    `json:"latitude"`
	Longitude   float64    `json:"longitude"`
	DepthMeters float64    `json:"depth_meters"`
	Status      BuoyStatus `json:"status"`
	LastSeen    time.Time  `json:"last_seen"`
	Notes       string     `json:"notes"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Sensor struct {
	ID          ID           `json:"id"`
	BuoyID      ID           `json:"buoy_id"`
	Serial      string       `json:"serial"`
	Kind        SensorKind   `json:"kind"`
	Status      SensorStatus `json:"status"`
	SampleRate  float64      `json:"sample_rate"`
	Calibration float64      `json:"calibration"`
	LastReading time.Time    `json:"last_reading"`
	CreatedAt   time.Time    `json:"created_at"`
}

type CalibrationBatch struct {
	ID          ID          `json:"id"`
	BuoyID      ID          `json:"buoy_id"`
	ProtocolID  ID          `json:"protocol_id"`
	Status      BatchStatus `json:"status"`
	WindowStart time.Time   `json:"window_start"`
	WindowEnd   time.Time   `json:"window_end"`
	StartedAt   time.Time   `json:"started_at"`
	CompletedAt time.Time   `json:"completed_at"`
	Reviewer    string      `json:"reviewer"`
	Summary     string      `json:"summary"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type AcousticSample struct {
	ID          ID        `json:"id"`
	BatchID     ID        `json:"batch_id"`
	SensorID    ID        `json:"sensor_id"`
	CapturedAt  time.Time `json:"captured_at"`
	FrequencyHz float64   `json:"frequency_hz"`
	AmplitudeDB float64   `json:"amplitude_db"`
	NoiseDB     float64   `json:"noise_db"`
	DurationMS  int       `json:"duration_ms"`
	Sequence    int       `json:"sequence"`
	PayloadHash string    `json:"payload_hash"`
	Valid       bool      `json:"valid"`
}

type QualityResult struct {
	ID           ID        `json:"id"`
	BatchID      ID        `json:"batch_id"`
	Coverage     float64   `json:"coverage"`
	NoiseFloorDB float64   `json:"noise_floor_db"`
	DriftDB      float64   `json:"drift_db"`
	Continuity   float64   `json:"continuity"`
	Score        float64   `json:"score"`
	Decision     Decision  `json:"decision"`
	Reasons      []string  `json:"reasons"`
	EvaluatedAt  time.Time `json:"evaluated_at"`
	Evaluator    string    `json:"evaluator"`
}

type Alert struct {
	ID             ID         `json:"id"`
	BuoyID         ID         `json:"buoy_id"`
	BatchID        ID         `json:"batch_id"`
	Severity       Severity   `json:"severity"`
	State          AlertState `json:"state"`
	Code           string     `json:"code"`
	Message        string     `json:"message"`
	OpenedAt       time.Time  `json:"opened_at"`
	AcknowledgedAt time.Time  `json:"acknowledged_at"`
	ClosedAt       time.Time  `json:"closed_at"`
	Owner          string     `json:"owner"`
}

type Protocol struct {
	ID             ID            `json:"id"`
	Name           string        `json:"name"`
	Version        int           `json:"version"`
	MinFrequencyHz float64       `json:"min_frequency_hz"`
	MaxFrequencyHz float64       `json:"max_frequency_hz"`
	MinDurationMS  int           `json:"min_duration_ms"`
	MaxDurationMS  int           `json:"max_duration_ms"`
	WindowMinutes  int           `json:"window_minutes"`
	State          ProtocolState `json:"state"`
	CreatedAt      time.Time     `json:"created_at"`
}

type ExportJob struct {
	ID          ID           `json:"id"`
	BatchID     ID           `json:"batch_id"`
	Format      ExportFormat `json:"format"`
	State       ExportState  `json:"state"`
	RequestedBy string       `json:"requested_by"`
	Path        string       `json:"path"`
	CreatedAt   time.Time    `json:"created_at"`
	FinishedAt  time.Time    `json:"finished_at"`
}

type AuditEvent struct {
	ID        ID        `json:"id"`
	Entity    string    `json:"entity"`
	EntityID  ID        `json:"entity_id"`
	Action    string    `json:"action"`
	Actor     string    `json:"actor"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}
