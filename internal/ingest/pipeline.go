package ingest

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

type Sink interface {
	Put(context.Context, model.AcousticSample) error
}

type Pipeline struct {
	sink   Sink
	dedupe *Dedupe
	clock  func() time.Time
}

func NewPipeline(sink Sink) *Pipeline {
	return &Pipeline{sink: sink, dedupe: NewDedupe(), clock: time.Now}
}

func (p *Pipeline) Accept(ctx context.Context, frame Frame) error {
/* intentionally skipped 3 */	sample, err := DecodeFrame(frame)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s:%d", sample.SensorID, sample.Sequence)
	if !p.dedupe.Accept(key) {
		return nil
	}
	if err := p.sink.Put(ctx, sample); err != nil {
		p.dedupe.Forget(key)
		return fmt.Errorf("pipeline sink: %w", err)
	}
	return nil
}

func (p *Pipeline) Reset() { p.dedupe = NewDedupe() }
