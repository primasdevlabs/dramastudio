package services

import (
	"context"
	"fmt"
	"time"

	"dramastudio/internal/continuity/domain"
)

type ContinuityService struct {
	repo domain.ContinuityRepository
}

func NewContinuityService(repo domain.ContinuityRepository) *ContinuityService {
	return &ContinuityService{repo: repo}
}

func (s *ContinuityService) RunCheck(ctx context.Context, projectID, episodeID, sceneID, category string, severity domain.Severity, entity, expected, actual, cause, evidence string) (*domain.ContinuityIssue, error) {
	id := fmt.Sprintf("issue_%d", time.Now().UnixNano())
	issue := &domain.ContinuityIssue{
		ID:            id,
		ProjectID:     projectID,
		EpisodeID:     episodeID,
		SceneID:       sceneID,
		Category:      category,
		Severity:      severity,
		Entity:        entity,
		ExpectedState: expected,
		ActualState:   actual,
		Cause:         cause,
		Evidence:      evidence,
		IsResolved:    false,
	}
	if err := s.repo.SaveIssue(ctx, issue); err != nil {
		return nil, err
	}
	return issue, nil
}

func (s *ContinuityService) ListIssues(ctx context.Context, projectID string) ([]*domain.ContinuityIssue, error) {
	return s.repo.ListIssuesByProject(ctx, projectID)
}

func (s *ContinuityService) ResolveIssue(ctx context.Context, id, resolution string) error {
	return s.repo.ResolveIssue(ctx, id, resolution)
}
