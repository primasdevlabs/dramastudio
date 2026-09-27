package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// LocalStorage writes objects under a root directory on disk. Intended for
// local development; SignedURL returns a file-scoped pseudo URL.
type LocalStorage struct {
	root string
}

func NewLocal(root string) (*LocalStorage, error) {
	if root == "" {
		return nil, fmt.Errorf("storage: local root dir required")
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	return &LocalStorage{root: root}, nil
}

func (s *LocalStorage) path(key string) string {
	return filepath.Join(s.root, filepath.FromSlash(key))
}

func (s *LocalStorage) Upload(_ context.Context, key string, r io.Reader, _ int64, _ string) error {
	p := s.path(key)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}

func (s *LocalStorage) Download(_ context.Context, key string) (io.ReadCloser, error) {
	return os.Open(s.path(key))
}

func (s *LocalStorage) Delete(_ context.Context, key string) error {
	return os.Remove(s.path(key))
}

func (s *LocalStorage) SignedURL(_ context.Context, key string, _ time.Duration) (string, error) {
	return (&url.URL{Scheme: "file", Path: s.path(key)}).String(), nil
}
