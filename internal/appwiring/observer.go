package appwiring

import (
	"context"
	"os"

	"dramastudio/internal/agents/director"
	contsvc "dramastudio/internal/continuity/application/services"
	continuitydomain "dramastudio/internal/continuity/domain"
	prodsvc "dramastudio/internal/production/application/services"
	proddomain "dramastudio/internal/production/domain"
)

// ProductionObserver implements director.Observer by reading live
// production and continuity state — the "observe" step of the director loop.
type ProductionObserver struct {
	production *prodsvc.ProductionService
	continuity *contsvc.ContinuityService
}

func NewProductionObserver(p *prodsvc.ProductionService, c *contsvc.ContinuityService) *ProductionObserver {
	return &ProductionObserver{production: p, continuity: c}
}

func (o *ProductionObserver) Observe(ctx context.Context, projectID, episodeID string) (*director.Observation, error) {
	obs := &director.Observation{ProjectID: projectID, EpisodeID: episodeID}
	if err := o.observeRun(ctx, projectID, episodeID, obs); err != nil {
		return nil, err
	}
	if err := o.observeApprovals(ctx, projectID, obs); err != nil {
		return nil, err
	}
	if err := o.observeIssues(ctx, projectID, episodeID, obs); err != nil {
		return nil, err
	}
	if err := o.observeJobs(ctx, projectID, episodeID, obs); err != nil {
		return nil, err
	}
	return obs, nil
}

// observeRun records the status of the episode's latest production run.
func (o *ProductionObserver) observeRun(ctx context.Context, projectID, episodeID string, obs *director.Observation) error {
	runs, err := o.production.ListRuns(ctx, projectID)
	if err != nil {
		return err
	}
	var latest *proddomain.ProductionRun
	for _, run := range runs {
		if run.EpisodeID != episodeID {
			continue
		}
		if latest == nil || run.CreatedAt.After(latest.CreatedAt) {
			latest = run
		}
	}
	if latest != nil {
		obs.RunStatus = string(latest.Status)
	}
	return nil
}

func (o *ProductionObserver) observeApprovals(ctx context.Context, projectID string, obs *director.Observation) error {
	approvals, err := o.production.ListApprovals(ctx, projectID, true)
	if err != nil {
		return err
	}
	obs.PendingApprovals = len(approvals)
	return nil
}

func (o *ProductionObserver) observeIssues(ctx context.Context, projectID, episodeID string, obs *director.Observation) error {
	issues, err := o.continuity.ListIssues(ctx, projectID, continuitydomain.IssueFilter{
		EpisodeID: episodeID,
		Status:    continuitydomain.IssueOpen,
	})
	if err != nil {
		return err
	}
	for _, i := range issues {
		obs.OpenIssues++
		if i.Severity == continuitydomain.SeverityBlocking {
			obs.BlockingIssues++
		}
	}
	return nil
}

func (o *ProductionObserver) observeJobs(ctx context.Context, projectID, episodeID string, obs *director.Observation) error {
	jobs, err := o.production.ListJobs(ctx, projectID)
	if err != nil {
		return err
	}
	for _, j := range jobs {
		if j.EpisodeID == episodeID && j.Status == proddomain.JobStatusFailed {
			obs.FailedJobs++
		}
	}
	return nil
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
