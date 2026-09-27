package domain

import "context"

type PublishingRepository interface {
	SavePublication(ctx context.Context, pub *Publication) error
}
