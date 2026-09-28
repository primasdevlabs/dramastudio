package domain

import (
	"context"
	"errors"
)

var ErrExecutionNotFound = errors.New("intelligence execution record not found")

// ExecutionRepository persists execution audit records.
type ExecutionRepository interface {
	SaveExecution(ctx context.Context, r *ExecutionRecord) error
	FindExecution(ctx context.Context, id string) (*ExecutionRecord, error)
	ListExecutions(ctx context.Context, projectID string) ([]*ExecutionRecord, error)
}

// PolicyLayerRepository persists operator-added policy layers (project,
// episode, task overrides) on top of the catalog's system/studio defaults.
type PolicyLayerRepository interface {
	SavePolicyLayer(ctx context.Context, p *Policy) error
	// ListPolicyLayers returns all layers for a policy ID at or below the
	// given scope ("episode:ep1" matches layer=episode scope_id=ep1 plus
	// wildcard project/studio/system layers when scopeID prefixes match).
	ListPolicyLayers(ctx context.Context, policyID string) ([]*Policy, error)
}
