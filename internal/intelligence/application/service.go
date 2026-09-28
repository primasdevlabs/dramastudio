// Package application exposes the production-intelligence subsystem:
// definition listing, execution-context preview, task execution, and
// operator policy-layer management.
package application

import (
	"context"
	"encoding/json"

	aiapp "dramastudio/internal/ai/application"
	aidomain "dramastudio/internal/ai/domain"
	"dramastudio/internal/intelligence/domain"
	"dramastudio/internal/intelligence/execution"
	"dramastudio/internal/intelligence/resolver"
)

// aiGenerator adapts the AI ExecuteGenerationHandler to execution.Generator
// so model routing, fallback, and job provenance stay in one place.
type aiGenerator struct {
	h *aiapp.ExecuteGenerationHandler
}

func NewAIGenerator(h *aiapp.ExecuteGenerationHandler) execution.Generator {
	return &aiGenerator{h: h}
}

func (g *aiGenerator) Generate(ctx context.Context, task domain.TaskContext, prompt string) (*execution.GenResult, error) {
	job, err := g.h.Handle(ctx, aiapp.ExecuteGenerationCommand{
		ProjectID:  task.ProjectID,
		Capability: aidomain.AICapability(task.Capability),
		Scope:      task.Scope,
		Input:      task.Input,
		Prompt:     prompt,
	})
	if err != nil {
		return nil, err
	}
	return &execution.GenResult{
		JobID: job.ID, ProviderID: job.ProviderID, ModelID: string(job.ModelID),
		ModelVersion: job.ModelVersion, Output: json.RawMessage(job.Output),
		Cost: job.Cost, Status: string(job.Status),
	}, nil
}

// IntelligenceService is the application facade.
type IntelligenceService struct {
	reg        *resolver.Registry
	exec       *execution.Executor
	executions domain.ExecutionRepository
	layers     domain.PolicyLayerRepository
}

func NewIntelligenceService(reg *resolver.Registry, exec *execution.Executor,
	executions domain.ExecutionRepository, layers domain.PolicyLayerRepository) *IntelligenceService {
	return &IntelligenceService{reg: reg, exec: exec, executions: executions, layers: layers}
}

func (s *IntelligenceService) Agents() []*domain.Agent         { return s.reg.Agents() }
func (s *IntelligenceService) Skills() []*domain.Skill         { return s.reg.Skills() }
func (s *IntelligenceService) Rules() []*domain.Rule           { return s.reg.Rules() }
func (s *IntelligenceService) Evaluators() []*domain.Evaluator { return s.reg.Evaluators() }
func (s *IntelligenceService) Policies() map[string][]*domain.Policy {
	return s.reg.Policies()
}

// Preview returns the resolved execution context for inspection.
func (s *IntelligenceService) Preview(ctx context.Context, agentID string, task domain.TaskContext) (*domain.AgentExecutionContext, error) {
	return s.exec.Preview(ctx, agentID, task)
}

// Execute runs a task through the intelligence pipeline.
func (s *IntelligenceService) Execute(ctx context.Context, agentID string, task domain.TaskContext) (*domain.ExecutionRecord, *execution.GenResult, error) {
	return s.exec.Run(ctx, agentID, task)
}

// SavePolicyLayer adds an operator override layer to a catalog policy.
func (s *IntelligenceService) SavePolicyLayer(ctx context.Context, p *domain.Policy) error {
	return s.layers.SavePolicyLayer(ctx, p)
}

func (s *IntelligenceService) ListExecutions(ctx context.Context, projectID string) ([]*domain.ExecutionRecord, error) {
	return s.executions.ListExecutions(ctx, projectID)
}
