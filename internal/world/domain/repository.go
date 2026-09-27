package domain

import "context"

type WorldRepository interface {
	SaveLocation(ctx context.Context, loc *Location) error
	FindLocationByID(ctx context.Context, id string) (*Location, error)
	ListLocations(ctx context.Context, projectID string) ([]*Location, error)
}
