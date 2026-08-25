package api

import "time"

type CreateBuoyRequest struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	DepthMeters float64 `json:"depth_meters"`
	Notes       string  `json:"notes"`
}

type CreateSensorRequest struct {
	ID         string  `json:"id"`
	BuoyID     string  `json:"buoy_id"`
	Serial     string  `json:"serial"`
	Kind       string  `json:"kind"`
	SampleRate float64 `json:"sample_rate"`
}

type CreateProtocolRequest struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Version        int     `json:"version"`
	MinFrequencyHz float64 `json:"min_frequency_hz"`
	MaxFrequencyHz float64 `json:"max_frequency_hz"`
	MinDurationMS  int     `json:"min_duration_ms"`
	MaxDurationMS  int     `json:"max_duration_ms"`
	WindowMinutes  int     `json:"window_minutes"`
}

type CreateBatchRequest struct {
	ID          string    `json:"id"`
	BuoyID      string    `json:"buoy_id"`
	ProtocolID  string    `json:"protocol_id"`
	WindowStart time.Time `json:"window_start"`
	WindowEnd   time.Time `json:"window_end"`
}

type SampleRequest struct {
	ID          string    `json:"id"`
	SensorID    string    `json:"sensor_id"`
	CapturedAt  time.Time `json:"captured_at"`
	FrequencyHz float64   `json:"frequency_hz"`
	AmplitudeDB float64   `json:"amplitude_db"`
	NoiseDB     float64   `json:"noise_db"`
	DurationMS  int       `json:"duration_ms"`
	Sequence    int       `json:"sequence"`
	PayloadHash string    `json:"payload_hash"`
	Valid       bool      `json:"valid"`
}

type CloseAlertRequest struct {
	Owner string `json:"owner"`
}

type ExportRequest struct {
	Format      string `json:"format"`
	RequestedBy string `json:"requested_by"`
}
