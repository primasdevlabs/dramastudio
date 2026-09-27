package domain

import "context"

type AgentRepository interface {
	FindByID(ctx context.Context, id AgentID) (*Agent, error)
	Save(ctx context.Context, agent *Agent) error
}
