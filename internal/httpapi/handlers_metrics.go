package httpapi

import (
	"net/http"

	"github.com/jb843051627/fjord-resonance/internal/api"
)

func (s *Server) metrics(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		method(writer, request, http.MethodGet)
		return
	}
	snapshot, err := s.app.Store.HealthSnapshot(request.Context())
	if err != nil {
		api.WriteError(writer, http.StatusInternalServerError, err)
		return
	}
	_ = api.WriteJSON(writer, http.StatusOK, snapshot)
}
