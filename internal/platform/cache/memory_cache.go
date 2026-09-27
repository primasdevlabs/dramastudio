package cache

import (
	"context"
	"sync"
	"time"
)

type memoryEntry struct {
	value     []byte
	expiresAt time.Time
}

// MemoryCache is an in-process cache for development and tests.
type MemoryCache struct {
	mu      sync.RWMutex
	entries map[string]memoryEntry
}

func NewMemory() *MemoryCache {
	return &MemoryCache{entries: make(map[string]memoryEntry)}
}

func (c *MemoryCache) Get(_ context.Context, key string) ([]byte, error) {
	c.mu.RLock()
	e, ok := c.entries[key]
	c.mu.RUnlock()
	if !ok || (!e.expiresAt.IsZero() && time.Now().After(e.expiresAt)) {
		return nil, ErrMiss
	}
	return e.value, nil
}

func (c *MemoryCache) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	e := memoryEntry{value: value}
	if ttl > 0 {
		e.expiresAt = time.Now().Add(ttl)
	}
	c.mu.Lock()
	c.entries[key] = e
	c.mu.Unlock()
	return nil
}

func (c *MemoryCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	delete(c.entries, key)
	c.mu.Unlock()
	return nil
}
