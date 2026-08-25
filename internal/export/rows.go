package export

import (
	"fmt"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

type SampleRow struct {
	BatchID    string
	SensorID   string
	CapturedAt time.Time
	Sequence   int
	Frequency  float64
	Amplitude  float64
	Noise      float64
	Valid      bool
}

func NewRows(samples []model.AcousticSample) []SampleRow {
	rows := make([]SampleRow, 0, len(samples))
	for _, sample := range samples {
		rows = append(rows, SampleRow{BatchID: string(sample.BatchID), SensorID: string(sample.SensorID), CapturedAt: sample.CapturedAt, Sequence: sample.Sequence, Frequency: sample.FrequencyHz, Amplitude: sample.AmplitudeDB, Noise: sample.NoiseDB, Valid: sample.Valid})
	}
	return rows
}

func FormatFloat(value float64) string { return fmt.Sprintf("%.3f", value) }
