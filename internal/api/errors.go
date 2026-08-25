package api

import (
	"errors"
	"net/http"

	"github.com/jb843051627/fjord-resonance/internal/store"
)

func StatusForError(err error) int {
	switch {
	case errors.Is(err, store.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, store.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, store.ErrInvalidState), err == store.ErrValidation:
		return http.StatusUnprocessableEntity
	case errors.Is(err, store.ErrCancelled):
		return http.StatusRequestTimeout
	default:
		return http.StatusInternalServerError
	}
}
