package services

import (
	"context"
	"dramastudio/internal/projects/domain"
)

type ProjectService struct {
	repo domain.ProjectRepository
}

func NewProjectService(repo domain.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) GetProject(ctx context.Context, id domain.ProjectID) (*domain.Project, error) {
	return s.repo.FindByID(ctx, id)
}
