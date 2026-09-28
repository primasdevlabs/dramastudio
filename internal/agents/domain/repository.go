package domain

import (
	"context"
	"errors"
)

var (
	ErrDefinitionNotFound = errors.New("agent definition not found")
	ErrTaskNotFound       = errors.New("agent task not found")
	ErrInvalidTransition  = errors.New("invalid task state transition")
)

type AgentRepository interface {
	SaveDefinition(ctx context.Context, d *Definition) error
	FindDefinitionByID(ctx context.Context, id string) (*Definition, error)
	ListDefinitions(ctx context.Context) ([]*Definition, error)

	SaveTask(ctx context.Context, t *Task) error
	FindTaskByID(ctx context.Context, id string) (*Task, error)
	ListTasks(ctx context.Context, projectID string, status TaskStatus) ([]*Task, error)

	SaveDecision(ctx context.Context, d *Decision) error
	ListDecisionsByProject(ctx context.Context, projectID string) ([]*Decision, error)
}
