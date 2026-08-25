package model

import (
	"fmt"
	"sort"
)

type ProtocolWindow struct {
	OffsetMinutes int
	DurationMS    int
	FrequencyHz   float64
}

func (p Protocol) Window(step int) (ProtocolWindow, error) {
	if p.State != ProtocolReady {
		return ProtocolWindow{}, fmt.Errorf("protocol is not ready")
	}
	if step < 0 || step >= p.WindowMinutes {
		return ProtocolWindow{}, fmt.Errorf("protocol step is outside window")
	}
	span := p.MaxFrequencyHz - p.MinFrequencyHz
	frequency := p.MinFrequencyHz + span*(float64(step)+0.5)/float64(p.WindowMinutes)
	duration := p.MinDurationMS + (step % maxInt(p.MaxDurationMS-p.MinDurationMS+1, 1))
	return ProtocolWindow{OffsetMinutes: step, DurationMS: duration, FrequencyHz: frequency}, nil
}

func SortProtocols(protocols []Protocol) []Protocol {
	result := append([]Protocol(nil), protocols...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Name == result[j].Name {
			return result[i].Version > result[j].Version
		}
		return result[i].Name < result[j].Name
	})
	return result
}

func LatestProtocol(protocols []Protocol, name string) (Protocol, bool) {
	var selected Protocol
	found := false
	for _, protocol := range protocols {
		if protocol.Name == name && (!found || protocol.Version > selected.Version) {
			selected, found = protocol, true
		}
	}
	return selected, found
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
