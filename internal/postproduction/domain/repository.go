package domain

import "context"

type PostProductionRepository interface {
	FindEditByID(ctx context.Context, id string) (*Edit, error)
	SaveEdit(ctx context.Context, edit *Edit) error
}
