package ingest

import (
	"fmt"
	"sync"
)

type Checkpoint struct {
	mu      sync.RWMutex
	offsets map[string]int64
}

func NewCheckpoint() *Checkpoint { return &Checkpoint{offsets: make(map[string]int64)} }

func (c *Checkpoint) Save(stream string, offset int64) error {
	if stream == "" || offset < 0 {
		return fmt.Errorf("invalid checkpoint")
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.offsets[stream] = offset
	return nil
}

func (c *Checkpoint) Load(stream string) (int64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.offsets[stream]
	return value, ok
}

func (c *Checkpoint) Snapshot() map[string]int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make(map[string]int64, len(c.offsets))
	for key, value := range c.offsets {
		result[key] = value
	}
	return result
}
