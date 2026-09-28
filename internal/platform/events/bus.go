package events

import "context"

// EventBus is the domain event contract used by application services.
type EventBus interface {
	Publish(ctx context.Context, event DomainEvent) error
	Subscribe(eventType string, handler Handler)
	SubscribeAll(handler Handler)
}
