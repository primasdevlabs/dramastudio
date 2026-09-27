package domain

import "context"

type ModelRegistryRepository interface {
	FindModelByID(ctx context.Context, id ModelID) (*Model, error)
	FindPolicyByProject(ctx context.Context, projectID string) (*ModelPolicy, error)
	SaveModel(ctx context.Context, model *Model) error
	SavePolicy(ctx context.Context, policy *ModelPolicy) error
}
