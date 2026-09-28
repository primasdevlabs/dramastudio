package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrPublicationNotFound = errors.New("publication not found")
	ErrChannelNotFound     = errors.New("channel not found")
	ErrInvalidTransition   = errors.New("invalid publication state transition")
)

type PublishingRepository interface {
	SaveChannel(ctx context.Context, c *Channel) error
	FindChannelByID(ctx context.Context, id string) (*Channel, error)
	ListChannels(ctx context.Context, projectID string) ([]*Channel, error)

	SavePublication(ctx context.Context, pub *Publication) error
	FindPublicationByID(ctx context.Context, id string) (*Publication, error)
	FindPublicationByIdempotencyKey(ctx context.Context, key string) (*Publication, error)
	ListPublicationsByProject(ctx context.Context, projectID string) ([]*Publication, error)
	ListDueScheduled(ctx context.Context, now time.Time) ([]*Publication, error)
}
