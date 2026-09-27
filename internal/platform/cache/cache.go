package cache

import (
	"context"
	"time"
)

// ErrMiss is returned when a key does not exist in the cache.
var ErrMiss = errMiss{}

type errMiss struct{}

func (errMiss) Error() string { return "cache: miss" }

// Cache is a short-lived coordination/caching layer. It is never the
// source of truth for production state (see Backend.md §45).
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
