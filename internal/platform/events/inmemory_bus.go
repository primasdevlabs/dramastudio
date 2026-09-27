package events

import (
	"context"
	"sync"
)

// InMemoryBus dispatches domain events to in-process subscribers. Events are
// also persisted to the event log by the caller where durability matters.
type InMemoryBus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{handlers: make(map[string][]Handler)}
}

func (b *InMemoryBus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	handlers := append([]Handler(nil), b.handlers[event.EventType()]...)
	b.mu.RUnlock()
	for _, h := range handlers {
		if err := h(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

func (b *InMemoryBus) Subscribe(eventType string, handler Handler) {
	b.mu.Lock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
	b.mu.Unlock()
}
