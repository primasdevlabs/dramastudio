package services

import (
	"context"

	"dramastudio/internal/agents/director"
	"dramastudio/internal/agents/domain"
)

type AgentService struct {
	repo         domain.AgentRepository
	leadDirector *director.LeadDirector
}

func NewAgentService(repo domain.AgentRepository) *AgentService {
	return &AgentService{
		repo:         repo,
		leadDirector: director.NewLeadDirector(),
	}
}

func (s *AgentService) TriggerLeadDirectorStep(ctx context.Context, projectID, episodeID string) (*domain.Decision, error) {
	dec, err := s.leadDirector.ExecuteLoopStep(ctx, projectID, episodeID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.SaveDecision(ctx, dec); err != nil {
		return nil, err
	}
	return dec, nil
}

func (s *AgentService) ListDecisions(ctx context.Context, projectID string) ([]*domain.Decision, error) {
	return s.repo.ListDecisionsByProject(ctx, projectID)
}
