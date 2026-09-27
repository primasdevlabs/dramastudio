package application

import (
	"context"
	"dramastudio/internal/ai/domain"
)

type ResolveModelUseCase struct {
	repo domain.ModelRegistryRepository
}

func NewResolveModelUseCase(repo domain.ModelRegistryRepository) *ResolveModelUseCase {
	return &ResolveModelUseCase{repo: repo}
}

func (u *ResolveModelUseCase) Resolve(ctx context.Context, projectID string, cap domain.AICapability) (*domain.Model, error) {
	policy, err := u.repo.FindPolicyByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	modelID, ok := policy.Mapping[cap]
	if !ok {
		return nil, domain.ErrCapabilityNotSupported
	}
	return u.repo.FindModelByID(ctx, modelID)
}
