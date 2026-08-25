package httpapi

import (
	"net/http"

	"github.com/jb843051627/fjord-resonance/internal/api"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

func (s *Server) alertAction(writer http.ResponseWriter, request *http.Request) {
	parts := pathParts(request.URL.Path, "/api/v1/alerts/")
	if len(parts) != 2 || parts[1] != "close" || request.Method != http.MethodPost {
		api.WriteError(writer, http.StatusNotFound, store.ErrNotFound)
		return
	}
	var input api.CloseAlertRequest
	if err := api.DecodeJSON(request, &input); err != nil {
		api.WriteError(writer, http.StatusBadRequest, err)
		return
	}
	if err := s.app.Alerts.Close(request.Context(), modelID(parts[0]), input.Owner); err != nil {
		api.WriteError(writer, api.StatusForError(err), err)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}
