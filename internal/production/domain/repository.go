package domain

import (
	"context"
	"errors"
)

var (
	ErrRunNotFound       = errors.New("production run not found")
	ErrJobNotFound       = errors.New("production job not found")
	ErrShotNotFound      = errors.New("shot not found")
	ErrApprovalNotFound  = errors.New("approval request not found")
	ErrInvalidTransition = errors.New("invalid run state transition")
)

type ProductionRepository interface {
	// Productions
	SaveProduction(ctx context.Context, p *Production) error
	FindProductionByProject(ctx context.Context, projectID string) (*Production, error)

	// Runs
	SaveRun(ctx context.Context, run *ProductionRun) error
	FindRunByID(ctx context.Context, id string) (*ProductionRun, error)
	ListRunsByProject(ctx context.Context, projectID string) ([]*ProductionRun, error)

	// Jobs
	SaveJob(ctx context.Context, job *ProductionJob) error
	FindJobByID(ctx context.Context, id string) (*ProductionJob, error)
	FindJobByIdempotencyKey(ctx context.Context, key string) (*ProductionJob, error)
	ListJobsByProject(ctx context.Context, projectID string) ([]*ProductionJob, error)
	ListJobsByRun(ctx context.Context, runID string) ([]*ProductionJob, error)

	// Shots
	SaveShot(ctx context.Context, shot *Shot) error
	FindShotByID(ctx context.Context, id string) (*Shot, error)
	ListShots(ctx context.Context, episodeID, sceneID string) ([]*Shot, error)

	// Approvals
	SaveApproval(ctx context.Context, app *ApprovalRequest) error
	FindApprovalByID(ctx context.Context, id string) (*ApprovalRequest, error)
	ListApprovalsByProject(ctx context.Context, projectID string, pendingOnly bool) ([]*ApprovalRequest, error)
}
