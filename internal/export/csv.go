package export

import (
	"context"
	"encoding/csv"
	"io"
	"strconv"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

func ContextForExport(ctx context.Context) context.Context { return ctx }

func WriteSamples(w io.Writer, samples []model.AcousticSample) error {
	csvWriter := csv.NewWriter(w)
	if err := csvWriter.Write([]string{"batch_id", "sensor_id", "captured_at", "sequence", "frequency_hz", "amplitude_db", "noise_db", "valid"}); err != nil {
		return err
	}
	for _, sample := range samples {
		if err := csvWriter.Write([]string{string(sample.BatchID), string(sample.SensorID), sample.CapturedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"), strconv.Itoa(sample.Sequence), FormatFloat(sample.FrequencyHz), FormatFloat(sample.AmplitudeDB), FormatFloat(sample.NoiseDB), strconv.FormatBool(sample.Valid)}); err != nil {
			return err
		}
	}
	csvWriter.Flush()
	return csvWriter.Error()
}
