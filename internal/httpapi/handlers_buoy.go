package httpapi

import (
	"net/http"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/api"
	"github.com/jb843051627/fjord-resonance/internal/model"
)

func (s *Server) buoys(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		var input api.CreateBuoyRequest
		if err := api.DecodeJSON(request, &input); err != nil {
			api.WriteError(writer, http.StatusBadRequest, err)
			return
		}
		buoy, err := s.app.Buoys.Create(request.Context(), api.BuoyFromRequest(input, time.Now().UTC()))
		if err != nil {
			api.WriteError(writer, api.StatusForError(err), err)
			return
		}
		_ = api.WriteJSON(writer, http.StatusCreated, buoy)
	case http.MethodGet:
		items, err := s.app.Buoys.List(request.Context(), model.BuoyStatus(request.URL.Query().Get("status")), request.URL.Query().Get("name"), 100)
		if err != nil {
			api.WriteError(writer, api.StatusForError(err), err)
			return
		}
		_ = api.WriteJSON(writer, http.StatusOK, items)
	default:
		method(writer, request, http.MethodGet+", "+http.MethodPost)
	}
}
