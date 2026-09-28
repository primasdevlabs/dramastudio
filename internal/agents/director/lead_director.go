package director

import (
	"context"
	"time"

	"github.com/google/uuid"

	"dramastudio/internal/agents/domain"
)

// Observation is the production state the director inspects each loop
// iteration. Callers assemble it from the owning contexts (production,
// continuity, media) so agents never import sibling domain packages.
type Observation struct {
	ProjectID        string
	EpisodeID        string
	RunStatus        string // pending|running|paused|awaiting_approval|completed|failed|stopped
	PendingApprovals int
	BlockingIssues   int
	OpenIssues       int
	FailedJobs       int
	BudgetSpent      float64
	BudgetLimit      float64
}

// Decision verbs emitted by the loop.
const (
	DecideAdvanceStage      = "ADVANCE_PRODUCTION_STAGE"
	DecideWaitForApproval   = "WAIT_FOR_APPROVAL"
	DecideEscalate          = "ESCALATE_TO_HUMAN"
	DecideRetryFailed       = "RETRY_FAILED_JOBS"
	DecideResolveContinuity = "RESOLVE_CONTINUITY"
	DecideBudgetExceeded    = "HALT_BUDGET_EXCEEDED"
	DecideIdle              = "IDLE"
)

// Observer supplies the current production truth for a loop step.
type Observer interface {
	Observe(ctx context.Context, projectID, episodeID string) (*Observation, error)
}

// LeadDirector runs Observe → Assess → Decide, emitting an auditable
// Decision each step (§43). It decides; humans and workflows act.
type LeadDirector struct {
	observer Observer
	mode     string // monitored|autonomous
}

func NewLeadDirector(observer Observer, mode string) *LeadDirector {
	if mode == "" {
		mode = "monitored"
	}
	return &LeadDirector{observer: observer, mode: mode}
}

// ExecuteLoopStep observes production state and emits the next decision.
func (d *LeadDirector) ExecuteLoopStep(ctx context.Context, projectID, episodeID string) (*domain.Decision, error) {
	obs, err := d.observer.Observe(ctx, projectID, episodeID)
	if err != nil {
		return nil, err
	}
	return d.decide(obs), nil
}

// decide is the pure rule engine: highest-priority condition wins.
func (d *LeadDirector) decide(obs *Observation) *domain.Decision {
	verb, reason := assess(obs)
	return &domain.Decision{
		ID:            "dec_" + uuid.NewString(),
		ProjectID:     obs.ProjectID,
		EpisodeID:     obs.EpisodeID,
		Decision:      verb,
		Reason:        reason,
		DecisionMaker: domain.DecisionMakerLeadDirector,
		Mode:          d.mode,
		Timestamp:     time.Now().UTC(),
	}
}

func assess(obs *Observation) (verb, reason string) {
	switch {
	case obs.BudgetLimit > 0 && obs.BudgetSpent >= obs.BudgetLimit:
		return DecideBudgetExceeded, "production budget exhausted; halting pending human review"
	case obs.BlockingIssues > 0:
		return DecideEscalate, "blocking continuity issues require human resolution"
	case obs.PendingApprovals > 0:
		return DecideWaitForApproval, "approval gates pending; production holds"
	case obs.RunStatus == "paused" || obs.RunStatus == "stopped":
		return DecideIdle, "production run is paused or stopped"
	case obs.FailedJobs > 0:
		return DecideRetryFailed, "failed jobs detected; scheduling retry"
	case obs.OpenIssues > 0:
		return DecideResolveContinuity, "non-blocking continuity issues need agent attention"
	case obs.RunStatus == "completed":
		return DecideIdle, "episode production complete"
	default:
		return DecideAdvanceStage, "no blockers; advancing production stage"
	}
}
