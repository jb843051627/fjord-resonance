package regression

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "net/http/httptest"
    "path/filepath"
    "strings"
    "sync"
    "testing"
    "time"

    "github.com/jb843051627/fjord-resonance/internal/engine"
    "github.com/jb843051627/fjord-resonance/internal/httpapi"
    "github.com/jb843051627/fjord-resonance/internal/ingest"
    "github.com/jb843051627/fjord-resonance/internal/model"
    "github.com/jb843051627/fjord-resonance/internal/quality"
    "github.com/jb843051627/fjord-resonance/internal/service"
    "github.com/jb843051627/fjord-resonance/internal/sqlite"
    "github.com/jb843051627/fjord-resonance/internal/store"
)

func setup22(t *testing.T) (*service.Application, *sqlite.Store, model.Buoy, model.Protocol, model.Sensor, model.CalibrationBatch) {
    t.Helper()
    path := filepath.Join(t.TempDir(), "calibration.db")
    database, err := sqlite.Open(path)
    if err != nil { t.Fatal(err) }
    t.Cleanup(func() { _ = database.Close() })
    app := service.NewApplication(database)
    now := time.Date(2026, 8, 25, 8, 0, 0, 0, time.UTC)
    buoy := model.Buoy{ID: "buoy-22", Name: "North Fjord", Latitude: 64.1, Longitude: -21.9, DepthMeters: 180, Status: model.BuoyActive, CreatedAt: now, UpdatedAt: now}
    if _, err := app.Buoys.Create(context.Background(), buoy); err != nil { t.Fatal(err) }
    protocol := model.Protocol{ID: "protocol-22", Name: "winter hydrophone sweep", Version: 1, MinFrequencyHz: 100, MaxFrequencyHz: 1000, MinDurationMS: 100, MaxDurationMS: 2000, WindowMinutes: 15, State: model.ProtocolReady, CreatedAt: now}
    if err := database.CreateProtocol(context.Background(), protocol); err != nil { t.Fatal(err) }
    sensor := model.Sensor{ID: "sensor-22", BuoyID: buoy.ID, Serial: "HYDRO-22", Kind: model.SensorHydrophone, Status: model.SensorReady, SampleRate: 48000, Calibration: 1, CreatedAt: now}
    if _, err := app.Sensors.Create(context.Background(), sensor); err != nil { t.Fatal(err) }
    batch := model.CalibrationBatch{ID: "batch-22", BuoyID: buoy.ID, ProtocolID: protocol.ID, Status: model.BatchDraft, WindowStart: now, WindowEnd: now.Add(time.Hour), CreatedAt: now, UpdatedAt: now}
    if _, err := app.Batches.Create(context.Background(), batch); err != nil { t.Fatal(err) }
    if err := app.Batches.Queue(context.Background(), batch.ID); err != nil { t.Fatal(err) }
    batch.Status = model.BatchQueued
    return app, database, buoy, protocol, sensor, batch
}

func sample22(sensorID, batchID model.ID, sequence int, at time.Time, frequency float64) model.AcousticSample {
    return service.Sample(model.ID(fmt.Sprintf("sample-%s-%d", sensorID, sequence)), batchID, sensorID, sequence, at, frequency, -12, -40)
}

var (
    _ = json.Marshal
    _ = errors.Is
    _ = fmt.Sprintf
    _ = strings.NewReader
    _ = http.MethodGet
    _ = httptest.NewRecorder
    _ = sync.Once{}
    _ = engine.CanTransition
    _ = httpapi.New
    _ = ingest.NewDedupe
    _ = quality.DefaultThresholds
    _ = store.ErrValidation
)

func TestBug22_ReportAlertsAreIndependentSnapshots(t *testing.T) {
    app, database, _, _, _, batch := setup22(t); result := model.QualityResult{ID: "quality-22", BatchID: batch.ID, Decision: model.DecisionReject, Score: 0.1, EvaluatedAt: time.Now().UTC()}; if err := database.SaveQuality(context.Background(), result); err != nil { t.Fatal(err) }; alert, err := app.Alerts.OpenForQuality(context.Background(), batch, result); if err != nil { t.Fatal(err) }
    first, err := app.Reports.Snapshot(context.Background(), batch.ID); if err != nil { t.Fatal(err) }; first.Alerts[0].Severity = model.SeverityInfo; second, err := app.Reports.Snapshot(context.Background(), batch.ID); if err != nil { t.Fatal(err) }; if second.Alerts[0].Severity != alert.Severity { t.Fatal("report snapshot shares alert backing storage") }
}
