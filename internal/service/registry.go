package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jb843051627/fjord-resonance/internal/engine"
	"github.com/jb843051627/fjord-resonance/internal/model"
	"github.com/jb843051627/fjord-resonance/internal/quality"
	"github.com/jb843051627/fjord-resonance/internal/sqlite"
	"github.com/jb843051627/fjord-resonance/internal/store"
)

type Application struct {
	Store     *sqlite.Store
	Buoys     *BuoyService
	Sensors   *SensorService
	Batches   *BatchService
	Samples   *SampleService
	Quality   *QualityService
	Alerts    *AlertService
	Protocols *ProtocolService
	Exports   *ExportService
	Reports   *ReportService
	Health    *HealthService
	Workers   *WorkerService
}

func NewApplication(s *sqlite.Store) *Application {
	evaluator := quality.NewEvaluator(quality.DefaultThresholds(), nil)
	queue := engine.NewQueue(64)
	app := &Application{Store: s}
	app.Buoys = NewBuoyService(s)
	app.Sensors = NewSensorService(s)
	app.Protocols = NewProtocolService(s)
	app.Batches = NewBatchService(s, evaluator)
	app.Samples = NewSampleService(s)
	app.Quality = NewQualityService(s, evaluator)
	app.Alerts = NewAlertService(s)
	app.Exports = NewExportService(s)
	app.Reports = NewReportService(s)
	app.Health = NewHealthService(s)
	app.Workers = NewWorkerService(s, queue, app.Quality, app.Alerts)
	return app
}

func (a *Application) Close() error { return a.Store.Close() }

func notFound(err error, entity, id string) error {
	if errors.Is(err, store.ErrNotFound) {
		return store.NotFound(entity, id)
	}
	return err
}

func audit(ctx context.Context, repo *sqlite.Store, entity string, id model.ID, action, actor, details string) error {
	event := model.AuditEvent{ID: model.ID(fmt.Sprintf("audit-%s-%s-%s", entity, id, action)), Entity: entity, EntityID: id, Action: action, Actor: actor, Details: details}
	return repo.AppendAudit(ctx, event)
}

type guardedCache[T any] struct {
	mu     sync.RWMutex
	values map[model.ID]T
}

func newGuardedCache[T any]() guardedCache[T] { return guardedCache[T]{values: make(map[model.ID]T)} }

func (c *guardedCache[T]) put(id model.ID, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[id] = value
}

func (c *guardedCache[T]) get(id model.ID) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.values[id]
	return value, ok
}
