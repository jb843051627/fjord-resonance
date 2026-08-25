package httpapi

import (
	"net/http"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/api"
)

func (s *Server) health(writer http.ResponseWriter, request *http.Request) {
	if !method(writer, request, http.MethodGet) {
		return
	}
	_ = api.WriteJSON(writer, http.StatusOK, api.HealthResponse{Status: "ok", Time: time.Now().UTC().Format(time.RFC3339Nano)})
}
