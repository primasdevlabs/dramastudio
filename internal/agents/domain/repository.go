package domain

import "context"

type AgentRepository interface {
	SaveDecision(ctx context.Context, d *Decision) error
	ListDecisionsByProject(ctx context.Context, projectID string) ([]*Decision, error)
}
