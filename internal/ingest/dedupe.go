package ingest

import "sync"

type Dedupe struct {
	mu   sync.RWMutex
	seen map[string]struct{}
}

func NewDedupe() *Dedupe { return &Dedupe{seen: make(map[string]struct{})} }

func (d *Dedupe) Accept(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, exists := d.seen[key]; exists {
		return false
	}
	d.seen[key] = struct{}{}
	return true
}

func (d *Dedupe) Forget(key string) { d.mu.Lock(); defer d.mu.Unlock(); delete(d.seen, key) }

func (d *Dedupe) Size() int { d.mu.Lock(); defer d.mu.Unlock(); return len(d.seen) }
