package events

import "context"

// Handler consumes a published domain event.
type Handler func(ctx context.Context, event DomainEvent) error
