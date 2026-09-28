package events

import (
	"context"
	"log/slog"
	"sync"
)

// Bus dispatches domain events to in-process subscribers and appends every
// event to the durable log (§47, §48). Subscriber errors are returned to
// the caller; log-append failures are logged and never block dispatch.
type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
	all      []Handler
	store    EventLogStore
}

// NewBus builds a bus backed by the given durable event log; nil means
// dispatch-only (used by tests and the memory store profile).
func NewBus(store EventLogStore) *Bus {
	return &Bus{handlers: make(map[string][]Handler), store: store}
}

func (b *Bus) Publish(ctx context.Context, event DomainEvent) error {
	if b == nil {
		return nil
	}
	if b.store != nil {
		if err := b.store.Append(ctx, event); err != nil {
			slog.Warn("event log append failed", "type", event.Type, "err", err)
		}
	}
	b.mu.RLock()
	handlers := append([]Handler(nil), b.handlers[event.Type]...)
	handlers = append(handlers, b.all...)
	b.mu.RUnlock()
	for _, h := range handlers {
		if err := h(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bus) Subscribe(eventType string, handler Handler) {
	b.mu.Lock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
	b.mu.Unlock()
}

// SubscribeAll registers a wildcard handler for every event type — used by
// the SSE broker bridge (§58).
func (b *Bus) SubscribeAll(handler Handler) {
	b.mu.Lock()
	b.all = append(b.all, handler)
	b.mu.Unlock()
}

// Emit is the nil-safe helper services call at state transitions:
// `s.events.Emit(ctx, events.AssetApproved, ...)`.
func (b *Bus) Emit(ctx context.Context, eventType, projectID, aggregateID, message string) {
	if b == nil {
		return
	}
	_ = b.Publish(ctx, New(eventType, projectID, aggregateID, message))
}
