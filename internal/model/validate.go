package model

import (
	"strings"
	"time"
)

func (b Buoy) Validate() error {
	if strings.TrimSpace(string(b.ID)) == "" || strings.TrimSpace(b.Name) == "" {
		return ErrModelValidation("buoy id and name are required")
	}
	if b.DepthMeters <= 0 || b.DepthMeters > 12000 {
		return ErrModelValidation("buoy depth is outside the operating envelope")
	}
	if b.Latitude < -90 || b.Latitude > 90 || b.Longitude < -180 || b.Longitude > 180 {
		return ErrModelValidation("buoy coordinates are invalid")
	}
	return nil
}

func (s Sensor) Validate() error {
	if s.ID == "" || s.BuoyID == "" || strings.TrimSpace(s.Serial) == "" {
		return ErrModelValidation("sensor identity is incomplete")
	}
	if s.SampleRate <= 0 || s.SampleRate > 1000000 {
		return ErrModelValidation("sensor sample rate is invalid")
	}
	return nil
}

func (p Protocol) Validate() error {
	if p.ID == "" || strings.TrimSpace(p.Name) == "" || p.Version <= 0 {
		return ErrModelValidation("protocol identity is incomplete")
	}
	if p.MinFrequencyHz <= 0 || p.MaxFrequencyHz <= p.MinFrequencyHz {
		return ErrModelValidation("protocol frequency range is invalid")
	}
	if p.MinDurationMS <= 0 || p.MaxDurationMS < p.MinDurationMS || p.WindowMinutes <= 0 {
		return ErrModelValidation("protocol timing is invalid")
	}
	return nil
}

func (s AcousticSample) Validate() error {
	if s.ID == "" || s.BatchID == "" || s.SensorID == "" {
		return ErrModelValidation("sample identity is incomplete")
	}
	if s.CapturedAt.IsZero() || s.FrequencyHz <= 0 || s.DurationMS <= 0 {
		return ErrModelValidation("sample measurement is invalid")
	}
	return nil
}

func (b CalibrationBatch) WindowValid(now time.Time) bool {
	if b.WindowStart.IsZero() || b.WindowEnd.IsZero() || !b.WindowEnd.After(b.WindowStart) {
		return false
	}
	return !now.Before(b.WindowStart.Add(-15*time.Minute)) && !now.After(b.WindowEnd.Add(15*time.Minute))
}

type ErrModelValidation string

func (e ErrModelValidation) Error() string { return string(e) }
