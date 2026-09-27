package application

import (
	"context"
	"encoding/json"
	"dramastudio/internal/ai/domain"
)

type ExecuteGenerationCommand struct {
	ProjectID  string
	Capability domain.AICapability
	Payload    json.RawMessage
}

type ExecuteGenerationHandler struct {
	resolver *ResolveModelUseCase
}

func NewExecuteGenerationHandler(resolver *ResolveModelUseCase) *ExecuteGenerationHandler {
	return &ExecuteGenerationHandler{resolver: resolver}
}

func (h *ExecuteGenerationHandler) Handle(ctx context.Context, cmd ExecuteGenerationCommand) (*domain.GenerationJob, error) {
	model, err := h.resolver.Resolve(ctx, cmd.ProjectID, cmd.Capability)
	if err != nil {
		return nil, err
	}

	job := &domain.GenerationJob{
		Capability: cmd.Capability,
		ModelID:    model.ID,
		Payload:    cmd.Payload,
		Status:     domain.StatusCompleted,
		Result:     "sample_generation_output",
	}
	return job, nil
}
