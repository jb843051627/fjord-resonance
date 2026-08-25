package httpapi

import (
	"net/http"

	"github.com/jb843051627/fjord-resonance/internal/api"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

func (s *Server) notFound(writer http.ResponseWriter, request *http.Request) {
	api.WriteError(writer, http.StatusNotFound, store.ErrNotFound)
}

func (s *Server) internalError(writer http.ResponseWriter, err error) {
	api.WriteError(writer, http.StatusInternalServerError, err)
}
