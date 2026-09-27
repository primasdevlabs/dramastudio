package storage

import (
	"context"
	"io"
	"time"
)

// ObjectStorage stores large artifacts (images, video, audio, renders).
// PostgreSQL holds metadata only (Backend.md §44).
type ObjectStorage interface {
	Upload(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	SignedURL(ctx context.Context, key string, ttl time.Duration) (string, error)
}
