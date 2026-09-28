package persistence

import (
	"context"
	"sync"

	"dramastudio/internal/continuity/domain"
)

type InMemoryContinuityRepository struct {
	mu     sync.RWMutex
	checks map[string]*domain.ContinuityCheck
	issues map[string]*domain.ContinuityIssue
	events map[string]*domain.TimelineEvent
}

func NewInMemoryContinuityRepository() *InMemoryContinuityRepository {
	return &InMemoryContinuityRepository{
		checks: make(map[string]*domain.ContinuityCheck),
		issues: make(map[string]*domain.ContinuityIssue),
		events: make(map[string]*domain.TimelineEvent),
	}
}

func (r *InMemoryContinuityRepository) SaveCheck(ctx context.Context, c *domain.ContinuityCheck) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.checks[c.ID] = c
	return nil
}

func (r *InMemoryContinuityRepository) FindCheckByID(ctx context.Context, id string) (*domain.ContinuityCheck, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.checks[id]
	if !ok {
		return nil, domain.ErrCheckNotFound
	}
	return c, nil
}

func (r *InMemoryContinuityRepository) ListChecks(ctx context.Context, projectID, episodeID string) ([]*domain.ContinuityCheck, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.ContinuityCheck, 0)
	for _, c := range r.checks {
		if c.ProjectID != projectID {
			continue
		}
		if episodeID != "" && c.EpisodeID != episodeID {
			continue
		}
		res = append(res, c)
	}
	return res, nil
}

func (r *InMemoryContinuityRepository) SaveIssue(ctx context.Context, issue *domain.ContinuityIssue) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.issues[issue.ID] = issue
	return nil
}

func (r *InMemoryContinuityRepository) FindIssueByID(ctx context.Context, id string) (*domain.ContinuityIssue, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	i, ok := r.issues[id]
	if !ok {
		return nil, domain.ErrIssueNotFound
	}
	return i, nil
}

func (r *InMemoryContinuityRepository) ListIssues(ctx context.Context, projectID string, f domain.IssueFilter) ([]*domain.ContinuityIssue, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.ContinuityIssue, 0)
	for _, i := range r.issues {
		if i.ProjectID != projectID {
			continue
		}
		if f.EpisodeID != "" && i.EpisodeID != f.EpisodeID {
			continue
		}
		if f.Status != "" && i.Status != f.Status {
			continue
		}
		if f.Severity != "" && i.Severity != f.Severity {
			continue
		}
		if f.Category != "" && i.Category != f.Category {
			continue
		}
		res = append(res, i)
	}
	return res, nil
}

func (r *InMemoryContinuityRepository) SaveTimelineEvent(ctx context.Context, e *domain.TimelineEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events[e.ID] = e
	return nil
}

func (r *InMemoryContinuityRepository) ListTimelineEvents(ctx context.Context, projectID, episodeID string) ([]*domain.TimelineEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.TimelineEvent, 0)
	for _, e := range r.events {
		if e.ProjectID != projectID {
			continue
		}
		if episodeID != "" && e.EpisodeID != episodeID {
			continue
		}
		res = append(res, e)
	}
	return res, nil
}
