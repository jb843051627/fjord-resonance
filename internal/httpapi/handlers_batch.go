package httpapi

import (
	"net/http"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/api"
	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/service"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

func (s *Server) batches(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		items, err := s.app.Batches.List(request.Context(), store.BatchFilter{BuoyID: modelID(request.URL.Query().Get("buoy_id")), Limit: 100})
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
	var input api.CreateBatchRequest
	if err := api.DecodeJSON(request, &input); err != nil {
		api.WriteError(writer, http.StatusBadRequest, err)
		return
	}
	if input.WindowStart.IsZero() {
		input.WindowStart = time.Now().UTC()
	}
	if input.WindowEnd.IsZero() {
		input.WindowEnd = input.WindowStart.Add(time.Hour)
	}
	batch, err := s.app.Batches.Create(request.Context(), api.BatchFromRequest(input, time.Now().UTC()))
	if err != nil {
		api.WriteError(writer, http.StatusOK, err)
		return
	}
	_ = api.WriteJSON(writer, http.StatusCreated, batch)
}

func (s *Server) batchAction(writer http.ResponseWriter, request *http.Request) {
	parts := pathParts(request.URL.Path, "/api/v1/batches/")
	if len(parts) < 1 || parts[0] == "" {
		api.WriteError(writer, http.StatusNotFound, store.ErrNotFound)
		return
	}
	id := modelID(parts[0])
	if len(parts) == 1 && request.Method == http.MethodGet {
		batch, err := s.app.Batches.Get(request.Context(), id)
		if err != nil {
			api.WriteError(writer, api.StatusForError(err), err)
			return
		}
		_ = api.WriteJSON(writer, http.StatusOK, batch)
		return
	}
	if len(parts) != 2 || request.Method != http.MethodPost {
		api.WriteError(writer, http.StatusMethodNotAllowed, store.ErrInvalidState)
		return
	}
	var err error
	switch parts[1] {
	case "queue":
		err = s.app.Batches.Queue(request.Context(), id)
	case "start":
		err = s.app.Batches.Start(request.Context(), id, time.Now().UTC())
	case "release":
		err = s.app.Batches.Release(request.Context(), id, "api")
	case "evaluate":
		_, err = s.app.Batches.Evaluate(request.Context(), id, "api")
	case "samples":
		err = addSample(request, s.app, id)
	default:
		api.WriteError(writer, http.StatusNotFound, store.ErrNotFound)
		return
	}
	if err != nil {
		api.WriteError(writer, api.StatusForError(err), err)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

func addSample(request *http.Request, app *service.Application, batchID model.ID) error {
	var input api.SampleRequest
	if err := api.DecodeJSON(request, &input); err != nil {
		return err
	}
	return app.Samples.Add(request.Context(), api.SampleFromRequest(input, batchID))
}
