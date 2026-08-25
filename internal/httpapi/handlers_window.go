package httpapi

import (
	"net/http"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/api"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

func (s *Server) windowCheck(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		method(writer, request, http.MethodGet)
		return
	}
	from, _ := time.Parse(time.RFC3339, request.URL.Query().Get("from"))
	to, _ := time.Parse(time.RFC3339, request.URL.Query().Get("to"))
	if from.IsZero() || to.IsZero() {
		api.WriteError(writer, http.StatusBadRequest, store.ErrValidation)
		return
	}
	_ = api.WriteJSON(writer, http.StatusOK, map[string]any{"from": from, "to": to, "valid": to.After(from)})
}
