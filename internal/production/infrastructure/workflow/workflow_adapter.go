package workflow

import "context"

type ProductionWorkflowAdapter struct{}

func NewProductionWorkflowAdapter() *ProductionWorkflowAdapter {
	return &ProductionWorkflowAdapter{}
}

func (a *ProductionWorkflowAdapter) ExecuteRun(ctx context.Context, runID string) error {
	return nil
}
