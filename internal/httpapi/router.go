package httpapi

import (
	"net/http"
	"strings"

	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/service"
)

type Server struct {
	app *service.Application
	mux *http.ServeMux
}

func New(app *service.Application) *Server {
	s := &Server{app: app, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) Handler() http.Handler { return RequestID(s.mux) }

func (s *Server) routes() {
	s.mux.HandleFunc("/healthz", s.health)
	s.mux.HandleFunc("/api/v1/buoys", s.buoys)
	s.mux.HandleFunc("/api/v1/sensors", s.sensors)
	s.mux.HandleFunc("/api/v1/protocols", s.protocols)
	s.mux.HandleFunc("/api/v1/batches", s.batches)
	s.mux.HandleFunc("/api/v1/batches/", s.batchAction)
	s.mux.HandleFunc("/api/v1/alerts/", s.alertAction)
	s.mux.HandleFunc("/api/v1/exports/", s.exportAction)
}

func pathParts(path, prefix string) []string {
	return strings.Split(strings.Trim(strings.TrimPrefix(path, prefix), "/"), "/")
}

func modelID(value string) model.ID { return model.ID(strings.TrimSpace(value)) }
