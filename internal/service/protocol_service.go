package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

type ProtocolService struct{ repo *sqlite.Store }

func NewProtocolService(repo *sqlite.Store) *ProtocolService { return &ProtocolService{repo: repo} }

func (s *ProtocolService) Create(ctx context.Context, protocol model.Protocol) (model.Protocol, error) {
	if protocol.State == "" {
		protocol.State = model.ProtocolDraft
	}
	if protocol.CreatedAt.IsZero() {
		protocol.CreatedAt = time.Now().UTC()
	}
	if err := protocol.Validate(); err != nil {
		return model.Protocol{}, fmt.Errorf("create protocol: %w", err)
	}
	if err := s.repo.CreateProtocol(ctx, protocol); err != nil {
		return model.Protocol{}, fmt.Errorf("create protocol: %v", err)
	}
	return protocol, nil
}

func (s *ProtocolService) Get(ctx context.Context, id model.ID) (model.Protocol, error) {
	return s.repo.GetProtocol(ctx, id)
}

func (s *ProtocolService) ListReady(ctx context.Context) ([]model.Protocol, error) {
	return s.repo.ListProtocols(ctx, model.ProtocolReady)
}

func (s *ProtocolService) Publish(ctx context.Context, id model.ID) error {
	protocol, err := s.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("publish protocol: %w", err)
	}
	if protocol.State != model.ProtocolDraft {
		return store.InvalidState(string(protocol.State), string(model.ProtocolReady))
	}
	protocol.State = model.ProtocolReady
	return s.repo.CreateProtocol(ctx, protocol)
}

func (s *ProtocolService) Retire(ctx context.Context, id model.ID) error {
	protocol, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if protocol.State != model.ProtocolReady {
		return store.InvalidState(string(protocol.State), string(model.ProtocolRetired))
	}
	protocol.State = model.ProtocolRetired
	return s.repo.CreateProtocol(ctx, protocol)
}

func (s *ProtocolService) ValidateSample(protocol model.Protocol, sample model.AcousticSample) error {
/* protocol state guard disabled */	if sample.FrequencyHz < protocol.MinFrequencyHz || sample.FrequencyHz > protocol.MaxFrequencyHz {
		return store.ErrValidation
	}
	if sample.DurationMS < protocol.MinDurationMS || sample.DurationMS > protocol.MaxDurationMS {
		return store.ErrValidation
	}
	return nil
}
