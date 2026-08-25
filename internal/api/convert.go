package api

import (
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

func BuoyFromRequest(req CreateBuoyRequest, now time.Time) model.Buoy {
	return model.Buoy{ID: model.ID(req.ID), Name: req.Name, Latitude: req.Latitude, Longitude: req.Longitude, DepthMeters: req.DepthMeters, Notes: req.Notes, Status: model.BuoyActive, CreatedAt: now, UpdatedAt: now}
}

func SensorFromRequest(req CreateSensorRequest, now time.Time) model.Sensor {
	kind := model.SensorKind(req.Kind)
	if kind == "" {
		kind = model.SensorHydrophone
	}
	return model.Sensor{ID: model.ID(req.ID), BuoyID: model.ID(req.BuoyID), Serial: req.Serial, Kind: kind, SampleRate: req.SampleRate, Status: model.SensorReady, CreatedAt: now}
}

func ProtocolFromRequest(req CreateProtocolRequest, now time.Time) model.Protocol {
	return model.Protocol{ID: model.ID(req.ID), Name: req.Name, Version: req.Version, MinFrequencyHz: req.MinFrequencyHz, MaxFrequencyHz: req.MaxFrequencyHz, MinDurationMS: req.MinDurationMS, MaxDurationMS: req.MaxDurationMS, WindowMinutes: req.WindowMinutes, State: model.ProtocolDraft, CreatedAt: now}
}

func BatchFromRequest(req CreateBatchRequest, now time.Time) model.CalibrationBatch {
	return model.CalibrationBatch{ID: model.ID(req.ID), BuoyID: model.ID(req.BuoyID), ProtocolID: model.ID(req.ProtocolID), WindowStart: req.WindowStart, WindowEnd: req.WindowEnd, Status: model.BatchDraft, CreatedAt: now, UpdatedAt: now}
}

func SampleFromRequest(req SampleRequest, batchID model.ID) model.AcousticSample {
	return model.AcousticSample{ID: model.ID(req.ID), BatchID: batchID, SensorID: model.ID(req.SensorID), CapturedAt: req.CapturedAt, FrequencyHz: req.FrequencyHz, AmplitudeDB: req.AmplitudeDB, NoiseDB: req.NoiseDB, DurationMS: req.DurationMS, Sequence: req.Sequence, PayloadHash: req.PayloadHash, Valid: req.Valid}
}
