package messaging

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// RedisPubSub implements Publisher/Subscriber over Redis pub/sub. Used for
// realtime fan-out only; durable state stays in PostgreSQL (Backend.md §45).
type RedisPubSub struct {
	client *redis.Client
}

func NewRedisPubSub(client *redis.Client) *RedisPubSub {
	return &RedisPubSub{client: client}
}

func (p *RedisPubSub) Publish(ctx context.Context, topic string, msg []byte) error {
	return p.client.Publish(ctx, topic, msg).Err()
}

func (p *RedisPubSub) Subscribe(ctx context.Context, topic string) (<-chan []byte, error) {
	sub := p.client.Subscribe(ctx, topic)
	out := make(chan []byte, 16)
	go func() {
		defer close(out)
		defer sub.Close()
		ch := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case m, ok := <-ch:
				if !ok {
					return
				}
				select {
				case out <- []byte(m.Payload):
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out, nil
}
