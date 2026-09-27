package domain

import "context"

type ProductionRepository interface {
	FindByID(ctx context.Context, id string) (*Production, error)
	Save(ctx context.Context, prod *Production) error
}
