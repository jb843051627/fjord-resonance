package httpapi

import (
	"net/http"

	"github.com/jb843051627/fjord-resonance/internal/api"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

func (s *Server) reportAction(writer http.ResponseWriter, request *http.Request) {
	parts := pathParts(request.URL.Path, "/api/v1/reports/")
	if len(parts) != 1 || request.Method != http.MethodGet {
		api.WriteError(writer, http.StatusNotFound, store.ErrNotFound)
		return
	}
	report, err := s.app.Reports.Snapshot(request.Context(), modelID(parts[0]))
	if err != nil {
		api.WriteError(writer, api.StatusForError(err), err)
		return
	}
	_ = api.WriteJSON(writer, http.StatusOK, api.ReportResponse{Batch: report.Batch, Quality: report.Quality, Samples: report.Samples, Alerts: report.Alerts})
}
