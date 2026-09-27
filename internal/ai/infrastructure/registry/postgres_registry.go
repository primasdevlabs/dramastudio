package registry

import (
	"context"
	"dramastudio/internal/ai/domain"
)

type PostgresModelRegistry struct{}

func NewPostgresModelRegistry() *PostgresModelRegistry {
	return &PostgresModelRegistry{}
}

func (r *PostgresModelRegistry) FindModelByID(ctx context.Context, id domain.ModelID) (*domain.Model, error) {
	return &domain.Model{
		ID:         id,
		Name:       string(id),
		ProviderID: "provider_a",
		IsActive:   true,
	}, nil
}

func (r *PostgresModelRegistry) FindPolicyByProject(ctx context.Context, projectID string) (*domain.ModelPolicy, error) {
	policy := domain.NewModelPolicy(projectID)
	policy.Mapping[domain.CapabilityScript] = "script_model_v1"
	policy.Mapping[domain.CapabilityDialogue] = "dialogue_model_v1"
	policy.Mapping[domain.CapabilityStoryboard] = "storyboard_model_v1"
	policy.Mapping[domain.CapabilityAnimation] = "animation_model_v1"
	return policy, nil
}

func (r *PostgresModelRegistry) SaveModel(ctx context.Context, model *domain.Model) error {
	return nil
}

func (r *PostgresModelRegistry) SavePolicy(ctx context.Context, policy *domain.ModelPolicy) error {
	return nil
}
