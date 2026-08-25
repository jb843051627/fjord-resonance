package sqlite

import "context"

func (s *Store) initialize(ctx context.Context) error {
	statements := []string{
		`PRAGMA foreign_keys = ON`,
		`CREATE TABLE IF NOT EXISTS buoys (
            id TEXT PRIMARY KEY, name TEXT NOT NULL, latitude REAL NOT NULL, longitude REAL NOT NULL,
            depth_meters REAL NOT NULL, status TEXT NOT NULL, last_seen TEXT NOT NULL, notes TEXT NOT NULL,
            created_at TEXT NOT NULL, updated_at TEXT NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS sensors (
            id TEXT PRIMARY KEY, buoy_id TEXT NOT NULL REFERENCES buoys(id), serial TEXT NOT NULL UNIQUE,
            kind TEXT NOT NULL, status TEXT NOT NULL, sample_rate REAL NOT NULL, calibration REAL NOT NULL,
            last_reading TEXT NOT NULL, created_at TEXT NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS protocols (
            id TEXT PRIMARY KEY, name TEXT NOT NULL, version INTEGER NOT NULL, min_frequency REAL NOT NULL,
            max_frequency REAL NOT NULL, min_duration INTEGER NOT NULL, max_duration INTEGER NOT NULL,
            window_minutes INTEGER NOT NULL, state TEXT NOT NULL, created_at TEXT NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS batches (
            id TEXT PRIMARY KEY, buoy_id TEXT NOT NULL REFERENCES buoys(id), protocol_id TEXT NOT NULL REFERENCES protocols(id),
            status TEXT NOT NULL, window_start TEXT NOT NULL, window_end TEXT NOT NULL, started_at TEXT NOT NULL,
            completed_at TEXT NOT NULL, reviewer TEXT NOT NULL, summary TEXT NOT NULL, created_at TEXT NOT NULL, updated_at TEXT NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS samples (
            id TEXT PRIMARY KEY, batch_id TEXT NOT NULL REFERENCES batches(id), sensor_id TEXT NOT NULL REFERENCES sensors(id),
            captured_at TEXT NOT NULL, frequency REAL NOT NULL, amplitude REAL NOT NULL, noise REAL NOT NULL,
            duration_ms INTEGER NOT NULL, sequence_no INTEGER NOT NULL, payload_hash TEXT NOT NULL, valid INTEGER NOT NULL,
            UNIQUE(batch_id, sensor_id, sequence_no)
        )`,
		`CREATE TABLE IF NOT EXISTS quality_results (
            id TEXT PRIMARY KEY, batch_id TEXT NOT NULL UNIQUE REFERENCES batches(id), coverage REAL NOT NULL,
            noise_floor REAL NOT NULL, drift REAL NOT NULL, continuity REAL NOT NULL, score REAL NOT NULL,
            decision TEXT NOT NULL, reasons TEXT NOT NULL, evaluated_at TEXT NOT NULL, evaluator TEXT NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS alerts (
            id TEXT PRIMARY KEY, buoy_id TEXT NOT NULL REFERENCES buoys(id), batch_id TEXT NOT NULL,
            severity TEXT NOT NULL, state TEXT NOT NULL, code TEXT NOT NULL, message TEXT NOT NULL,
            opened_at TEXT NOT NULL, acknowledged_at TEXT NOT NULL, closed_at TEXT NOT NULL, owner TEXT NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS export_jobs (
            id TEXT PRIMARY KEY, batch_id TEXT NOT NULL REFERENCES batches(id), format TEXT NOT NULL,
            state TEXT NOT NULL, requested_by TEXT NOT NULL, path TEXT NOT NULL, created_at TEXT NOT NULL, finished_at TEXT NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS audit_events (
            id TEXT PRIMARY KEY, entity TEXT NOT NULL, entity_id TEXT NOT NULL, action TEXT NOT NULL,
            actor TEXT NOT NULL, details TEXT NOT NULL, created_at TEXT NOT NULL
        )`,
		`CREATE INDEX IF NOT EXISTS idx_batches_buoy_status ON batches(buoy_id, status)`,
		`CREATE INDEX IF NOT EXISTS idx_samples_batch_time ON samples(batch_id, captured_at)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_state ON alerts(state, severity)`,
	}
	for _, statement := range statements {
		if _, err := s.db.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}
