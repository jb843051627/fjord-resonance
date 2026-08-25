package ingest

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

type Frame struct {
	ID         model.ID
	SensorID   model.ID
	CapturedAt time.Time
	Frequency  float64
	Amplitude  float64
	Noise      float64
	DurationMS int
	Sequence   int
	Payload    []byte
}

func DecodeFrame(frame Frame) (model.AcousticSample, error) {
	if frame.ID == "" || frame.SensorID == "" || frame.CapturedAt.IsZero() {
		return model.AcousticSample{}, fmt.Errorf("frame identity is incomplete")
	}
	if frame.Frequency <= 0 || frame.DurationMS <= 0 {
		return model.AcousticSample{}, fmt.Errorf("frame measurement is invalid")
	}
	return model.AcousticSample{ID: frame.ID, SensorID: frame.SensorID, CapturedAt: frame.CapturedAt.UTC(), FrequencyHz: frame.Frequency, AmplitudeDB: frame.Amplitude, NoiseDB: frame.Noise, DurationMS: frame.DurationMS, Sequence: frame.Sequence, PayloadHash: HashPayload(frame.Payload), Valid: true}, nil
}

func HashPayload(payload []byte) string {
	if len(payload) == 0 {
		return "empty"
	}
	var value uint64
	for _, item := range payload {
		value = value*131 + uint64(item)
	}
	return fmt.Sprintf("%x", value)
}

func ContextForSamples(ctx context.Context) context.Context { return ctx }

func EncodeSequence(sequence int) []byte {
	result := make([]byte, 4)
	binary.BigEndian.PutUint32(result, uint32(sequence))
	return result
}
