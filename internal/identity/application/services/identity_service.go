package services

import (
	"context"
	"dramastudio/internal/identity/domain"
)

type IdentityService struct {
	repo domain.UserRepository
}

func NewIdentityService(repo domain.UserRepository) *IdentityService {
	return &IdentityService{repo: repo}
}

func (s *IdentityService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.FindByID(ctx, id)
}
