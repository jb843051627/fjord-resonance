package engine

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

type Lease struct {
	mu      sync.Mutex
	holders map[model.ID]string
	expiry  map[model.ID]time.Time
}

func NewLease() *Lease {
	return &Lease{holders: make(map[model.ID]string), expiry: make(map[model.ID]time.Time)}
}

func (l *Lease) Acquire(ctx context.Context, batchID model.ID, worker string, now time.Time, ttl time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if batchID == "" || worker == "" || ttl <= 0 {
		return fmt.Errorf("invalid lease")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if holder := l.holders[batchID]; holder != "" && l.expiry[batchID].After(now) && holder != worker {
		return fmt.Errorf("batch lease held by %s", holder)
	}
	l.holders[batchID], l.expiry[batchID] = worker, now.Add(ttl)
	return nil
}

func (l *Lease) Release(batchID model.ID, worker string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.holders[batchID] != worker {
		return false
	}
	delete(l.holders, batchID)
	delete(l.expiry, batchID)
	return true
}

func (l *Lease) Expired(batchID model.ID, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	expiry, ok := l.expiry[batchID]
	return !ok || !expiry.After(now)
}
