package store

import (
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

type BuoyFilter struct {
	Status model.BuoyStatus
	Name   string
	Limit  int
	Offset int
}

type BatchFilter struct {
	BuoyID model.ID
	Status model.BatchStatus
	From   time.Time
	To     time.Time
	Limit  int
	Offset int
}

type SampleFilter struct {
	From      time.Time
	To        time.Time
	OnlyValid bool
	Limit     int
	Offset    int
}

type AlertFilter struct {
	BuoyID   model.ID
	BatchID  model.ID
	State    model.AlertState
	Severity model.Severity
	Limit    int
	Offset   int
}

func NormalizeLimit(limit, fallback, max int) int {
	if limit <= 0 {
		limit = fallback
	}
	if limit > max {
		return max
	}
	return limit
}
