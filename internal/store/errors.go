package store

import "fmt"

var (
	ErrNotFound     = fmt.Errorf("fjord resonance: entity not found")
	ErrConflict     = fmt.Errorf("fjord resonance: entity conflict")
	ErrInvalidState = fmt.Errorf("fjord resonance: invalid state")
	ErrCancelled    = fmt.Errorf("fjord resonance: operation cancelled")
	ErrValidation   = fmt.Errorf("fjord resonance: validation failed")
)

func NotFound(entity, id string) error {
	return fmt.Errorf("%w: %s %s", ErrNotFound, entity, id)
}

func Conflict(entity, id string) error {
	return fmt.Errorf("%w: %s %s", ErrConflict, entity, id)
}

func InvalidState(from, to string) error {
	return fmt.Errorf("%w: %s -> %s", ErrInvalidState, from, to)
}
