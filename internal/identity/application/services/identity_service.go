package services

import (
	"context"
	"fmt"
	"time"

	"dramastudio/internal/identity/domain"
)

type IdentityService struct {
	repo domain.UserRepository
}

func NewIdentityService(repo domain.UserRepository) *IdentityService {
	return &IdentityService{repo: repo}
}

func (s *IdentityService) CreateUser(ctx context.Context, email, name, role, orgID string) (*domain.User, error) {
	id := fmt.Sprintf("user_%d", time.Now().UnixNano())
	if role == "" {
		role = "Producer"
	}
	u := &domain.User{
		ID:        id,
		Email:     email,
		Name:      name,
		Role:      role,
		OrgID:     orgID,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.repo.SaveUser(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *IdentityService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.FindUserByID(ctx, id)
}

func (s *IdentityService) ListUsers(ctx context.Context) ([]*domain.User, error) {
	return s.repo.ListUsers(ctx)
}
