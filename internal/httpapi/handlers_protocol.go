package httpapi

import (
	"net/http"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/api"
)

func (s *Server) protocols(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		items, err := s.app.Protocols.ListReady(request.Context())
		if err != nil {
			api.WriteError(writer, api.StatusForError(err), err)
			return
		}
		_ = api.WriteJSON(writer, http.StatusOK, items)
		return
	}
	if !method(writer, request, http.MethodPost) {
		return
	}
	var input api.CreateProtocolRequest
	if err := api.DecodeJSON(request, &input); err != nil {
		api.WriteError(writer, http.StatusBadRequest, err)
		return
	}
	protocol, err := s.app.Protocols.Create(request.Context(), api.ProtocolFromRequest(input, time.Now().UTC()))
	if err != nil {
		api.WriteError(writer, api.StatusForError(err), err)
		return
	}
	_ = api.WriteJSON(writer, http.StatusCreated, protocol)
}
