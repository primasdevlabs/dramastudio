package domain

import "context"

type ProductionRepository interface {
	SaveJob(ctx context.Context, job *ProductionJob) error
	FindJobByID(ctx context.Context, id string) (*ProductionJob, error)
	ListJobsByProject(ctx context.Context, projectID string) ([]*ProductionJob, error)
	SaveApproval(ctx context.Context, app *ApprovalRequest) error
	ListApprovalsByProject(ctx context.Context, projectID string) ([]*ApprovalRequest, error)
}
