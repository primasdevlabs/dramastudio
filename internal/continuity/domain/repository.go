package domain

import (
	"context"
	"errors"
)

var (
	ErrIssueNotFound = errors.New("continuity issue not found")
	ErrCheckNotFound = errors.New("continuity check not found")
	ErrEventNotFound = errors.New("timeline event not found")
)

type IssueFilter struct {
	EpisodeID string
	Status    IssueStatus
	Severity  Severity
	Category  string
}

type ContinuityRepository interface {
	SaveCheck(ctx context.Context, c *ContinuityCheck) error
	FindCheckByID(ctx context.Context, id string) (*ContinuityCheck, error)
	ListChecks(ctx context.Context, projectID, episodeID string) ([]*ContinuityCheck, error)

	SaveIssue(ctx context.Context, issue *ContinuityIssue) error
	FindIssueByID(ctx context.Context, id string) (*ContinuityIssue, error)
	ListIssues(ctx context.Context, projectID string, f IssueFilter) ([]*ContinuityIssue, error)

	SaveTimelineEvent(ctx context.Context, e *TimelineEvent) error
	ListTimelineEvents(ctx context.Context, projectID, episodeID string) ([]*TimelineEvent, error)
}
