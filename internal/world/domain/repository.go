package domain

import "context"

type WorldRepository interface {
	FindLocationByID(ctx context.Context, id string) (*Location, error)
	SaveLocation(ctx context.Context, location *Location) error
}
