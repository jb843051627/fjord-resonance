package store

import (
	"context"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

type BuoyRepository interface {
	CreateBuoy(context.Context, model.Buoy) error
	GetBuoy(context.Context, model.ID) (model.Buoy, error)
	ListBuoys(context.Context, BuoyFilter) ([]model.Buoy, error)
	UpdateBuoyStatus(context.Context, model.ID, model.BuoyStatus, time.Time) error
}

type SensorRepository interface {
	CreateSensor(context.Context, model.Sensor) error
	GetSensor(context.Context, model.ID) (model.Sensor, error)
	ListSensors(context.Context, model.ID) ([]model.Sensor, error)
	UpdateSensorCalibration(context.Context, model.ID, float64, model.SensorStatus) error
}

type ProtocolRepository interface {
	CreateProtocol(context.Context, model.Protocol) error
	GetProtocol(context.Context, model.ID) (model.Protocol, error)
	ListProtocols(context.Context, model.ProtocolState) ([]model.Protocol, error)
}

type BatchRepository interface {
	CreateBatch(context.Context, model.CalibrationBatch) error
	GetBatch(context.Context, model.ID) (model.CalibrationBatch, error)
	ListBatches(context.Context, BatchFilter) ([]model.CalibrationBatch, error)
	UpdateBatch(context.Context, model.CalibrationBatch) error
}

type SampleRepository interface {
	AddSample(context.Context, model.AcousticSample) error
	ListSamples(context.Context, model.ID, SampleFilter) ([]model.AcousticSample, error)
	CountSamples(context.Context, model.ID) (int, error)
}

type QualityRepository interface {
	SaveQuality(context.Context, model.QualityResult) error
	GetQuality(context.Context, model.ID) (model.QualityResult, error)
}

type AlertRepository interface {
	CreateAlert(context.Context, model.Alert) error
	GetAlert(context.Context, model.ID) (model.Alert, error)
	ListAlerts(context.Context, AlertFilter) ([]model.Alert, error)
	UpdateAlert(context.Context, model.Alert) error
}

type ExportRepository interface {
	CreateExport(context.Context, model.ExportJob) error
	GetExport(context.Context, model.ID) (model.ExportJob, error)
	UpdateExport(context.Context, model.ExportJob) error
}

type AuditRepository interface {
	AppendAudit(context.Context, model.AuditEvent) error
	ListAudit(context.Context, string, model.ID) ([]model.AuditEvent, error)
}

type Repository interface {
	BuoyRepository
	SensorRepository
	ProtocolRepository
	BatchRepository
	SampleRepository
	QualityRepository
	AlertRepository
	ExportRepository
	AuditRepository
}
