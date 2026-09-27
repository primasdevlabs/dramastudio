package domain

import "context"

type ContinuityRepository interface {
	SaveIssue(ctx context.Context, issue *ContinuityIssue) error
	FindIssueByID(ctx context.Context, id string) (*ContinuityIssue, error)
	ListIssuesByProject(ctx context.Context, projectID string) ([]*ContinuityIssue, error)
	ResolveIssue(ctx context.Context, id, resolution string) error
}
