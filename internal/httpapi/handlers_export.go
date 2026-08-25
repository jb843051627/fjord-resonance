package httpapi

import (
	"net/http"

	"github.com/jb843051627/fjord-resonance/internal/api"
	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

func (s *Server) exportAction(writer http.ResponseWriter, request *http.Request) {
	parts := pathParts(request.URL.Path, "/api/v1/exports/")
	if len(parts) != 2 || parts[1] != "csv" || request.Method != http.MethodGet {
		api.WriteError(writer, http.StatusNotFound, store.ErrNotFound)
		return
	}
	content, err := s.app.Exports.CSV(request.Context(), model.ID(parts[0]))
	if err != nil {
		api.WriteError(writer, api.StatusForError(err), err)
		return
	}
	writer.Header().Set("Content-Type", "text/csv; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(content)
}
