package events

import "context"

type EventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(eventType string, handler Handler)
}
