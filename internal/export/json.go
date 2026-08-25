package export

import (
	"encoding/json"
	"io"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

func WriteJSON(w io.Writer, value any) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

type JSONBundle struct {
	Batch   model.CalibrationBatch `json:"batch"`
	Quality model.QualityResult    `json:"quality"`
	Samples []model.AcousticSample `json:"samples"`
	Alerts  []model.Alert          `json:"alerts"`
}

func NewJSONBundle(batch model.CalibrationBatch, quality model.QualityResult, samples []model.AcousticSample, alerts []model.Alert) JSONBundle {
	return JSONBundle{Batch: batch, Quality: quality, Samples: model.CloneSamples(samples), Alerts: model.CloneAlerts(alerts)}
}
