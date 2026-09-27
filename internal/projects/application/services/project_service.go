package services

import (
	"context"
	"fmt"
	"time"

	"dramastudio/internal/projects/domain"
)

type ProjectService struct {
	repo domain.ProjectRepository
}

func NewProjectService(repo domain.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) CreateProject(ctx context.Context, name, description, genre, language string, mode domain.ProductionMode) (*domain.Project, error) {
	id := domain.ProjectID(fmt.Sprintf("proj_%d", time.Now().UnixNano()))
	project := domain.NewProject(id, name, description, genre, language, mode)
	if err := s.repo.Save(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) GetProject(ctx context.Context, id domain.ProjectID) (*domain.Project, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ProjectService) ListProjects(ctx context.Context) ([]*domain.Project, error) {
	return s.repo.ListAll(ctx)
}

func (s *ProjectService) SaveSeriesBible(ctx context.Context, projectID domain.ProjectID, premise, genre string, themes, worldRules, narrativeRules []string) (*domain.SeriesBible, error) {
	proj, err := s.repo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	version := 1
	if existing, err := s.repo.GetLatestBible(ctx, projectID); err == nil && existing != nil {
		version = existing.Version + 1
	}

	id := fmt.Sprintf("bible_%s_v%d", projectID, version)
	bible := domain.NewSeriesBible(id, proj.ID, version, premise)
	bible.Genre = genre
	bible.Themes = themes
	bible.WorldRules = worldRules
	bible.NarrativeRules = narrativeRules

	if err := s.repo.SaveBible(ctx, bible); err != nil {
		return nil, err
	}
	return bible, nil
}

func (s *ProjectService) GetLatestBible(ctx context.Context, projectID domain.ProjectID) (*domain.SeriesBible, error) {
	return s.repo.GetLatestBible(ctx, projectID)
}
