package engine

import (
	"sync"
	"time"
)

type Metrics struct {
	mu        sync.RWMutex
	completed int
	failed    int
	last      time.Time
}

func (m *Metrics) MarkCompleted(at time.Time) { m.mu.Lock(); m.completed++; m.last = at; m.mu.Unlock() }

func (m *Metrics) MarkFailed(at time.Time) { m.mu.Lock(); m.failed++; m.last = at; m.mu.Unlock() }

func (m *Metrics) Snapshot() (completed, failed int, last time.Time) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.completed, m.failed, m.last
}

func (m *Metrics) SuccessRate() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	total := m.completed + m.failed
	if total == 0 {
		return 0
	}
	return float64(m.completed) / float64(total)
}
