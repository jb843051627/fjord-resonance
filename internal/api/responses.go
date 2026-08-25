package api

import "github.com/jb843051627/fjord-resonance/internal/model"

type ErrorResponse struct {
	Error string `json:"error"`
}

type HealthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

type BatchResponse struct {
	Batch   model.CalibrationBatch `json:"batch"`
	Quality *model.QualityResult   `json:"quality,omitempty"`
}

type ReportResponse struct {
	Batch   model.CalibrationBatch `json:"batch"`
	Quality model.QualityResult    `json:"quality"`
	Samples []model.AcousticSample `json:"samples"`
	Alerts  []model.Alert          `json:"alerts"`
}
