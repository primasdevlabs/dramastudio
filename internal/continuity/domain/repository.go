package domain

import "context"

type ContinuityRepository interface {
	FindIssueByID(ctx context.Context, id string) (*ContinuityIssue, error)
	ListIssuesByEpisode(ctx context.Context, episodeID string) ([]*ContinuityIssue, error)
	SaveIssue(ctx context.Context, issue *ContinuityIssue) error
}
