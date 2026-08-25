package httpapi

import (
	"net/http"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/api"
)

func (s *Server) sensors(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		items, err := s.app.Sensors.List(request.Context(), modelID(request.URL.Query().Get("buoy_id")))
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
	var input api.CreateSensorRequest
	if err := api.DecodeJSON(request, &input); err != nil {
		api.WriteError(writer, http.StatusBadRequest, err)
		return
	}
	sensor, err := s.app.Sensors.Create(request.Context(), api.SensorFromRequest(input, time.Now().UTC()))
	if err != nil {
		api.WriteError(writer, api.StatusForError(err), err)
		return
	}
	_ = api.WriteJSON(writer, http.StatusCreated, sensor)
}
