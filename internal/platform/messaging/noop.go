package messaging

import "context"

// NoopPublisher drops messages; used when no message backend is configured.
type NoopPublisher struct{}

func (NoopPublisher) Publish(_ context.Context, _ string, _ []byte) error { return nil }
