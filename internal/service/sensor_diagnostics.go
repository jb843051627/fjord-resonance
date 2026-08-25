package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

type SensorDiagnostic struct {
	Sensor        model.Sensor
	LastHeartbeat time.Time
	Silent        bool
	Samples       int
	Detail        string
}

type SensorDiagnosticsService struct {
	repo    *sqlite.Store
	sensors *SensorService
}

func NewSensorDiagnosticsService(repo *sqlite.Store, sensors *SensorService) *SensorDiagnosticsService {
	return &SensorDiagnosticsService{repo: repo, sensors: sensors}
}

func (s *SensorDiagnosticsService) Inspect(ctx context.Context, id model.ID, now time.Time, silence time.Duration) (SensorDiagnostic, error) {
	sensor, err := s.sensors.Get(ctx, id)
	if err != nil {
		return SensorDiagnostic{}, fmt.Errorf("inspect sensor: %w", err)
	}
	heartbeat, ok := s.sensors.LastHeartbeat(id)
	if !ok {
		heartbeat = sensor.LastReading
	}
	count := 0
	if heartbeat.IsZero() {
		count = 0
	} else {
		count, err = s.countRecent(ctx, sensor.ID, now.Add(-silence), now)
		if err != nil {
			return SensorDiagnostic{}, err
		}
	}
	silent := heartbeat.IsZero() || now.Sub(heartbeat) > silence || count == 0
	detail := "sensor is producing samples"
	if silent {
		detail = "sensor has no recent acoustic sample"
	}
	return SensorDiagnostic{Sensor: sensor, LastHeartbeat: heartbeat, Silent: silent, Samples: count, Detail: detail}, nil
}

func (s *SensorDiagnosticsService) countRecent(ctx context.Context, sensorID model.ID, from, to time.Time) (int, error) {
	items, err := s.repo.ListSamples(ctx, "", store.SampleFilter{From: from, To: to, Limit: 10000})
	if err != nil {
		return 0, err
	}
	count := 0
	for _, item := range items {
		if item.SensorID == sensorID {
			count++
		}
	}
	return count, nil
}

func (s *SensorDiagnosticsService) Explain(diagnostic SensorDiagnostic) string {
	if diagnostic.Silent {
		return fmt.Sprintf("%s: %s", diagnostic.Sensor.ID, diagnostic.Detail)
	}
	return fmt.Sprintf("%s: %d recent samples", diagnostic.Sensor.ID, diagnostic.Samples)
}
