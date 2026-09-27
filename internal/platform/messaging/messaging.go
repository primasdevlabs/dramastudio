package messaging

import "context"

// Publisher emits messages to a topic for cross-process consumers.
type Publisher interface {
	Publish(ctx context.Context, topic string, msg []byte) error
}

// Subscriber receives messages from topics.
type Subscriber interface {
	Subscribe(ctx context.Context, topic string) (<-chan []byte, error)
}
