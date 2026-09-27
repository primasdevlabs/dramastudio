package domain

import "context"

type PublishingRepository interface {
	SavePublication(ctx context.Context, pub *Publication) error
	ListPublicationsByProject(ctx context.Context, projectID string) ([]*Publication, error)
}
