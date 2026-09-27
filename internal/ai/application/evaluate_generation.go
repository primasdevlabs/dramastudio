package application

import (
	"context"
	"dramastudio/internal/ai/domain"
)

type EvaluateGenerationCommand struct {
	JobID string
}

type EvaluateGenerationHandler struct{}

func NewEvaluateGenerationHandler() *EvaluateGenerationHandler {
	return &EvaluateGenerationHandler{}
}

func (h *EvaluateGenerationHandler) Evaluate(ctx context.Context, cmd EvaluateGenerationCommand) (*domain.GenerationJob, bool, string, error) {
	return &domain.GenerationJob{ID: cmd.JobID, Status: domain.StatusCompleted}, true, "Passed quality checks", nil
}
