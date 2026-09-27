package domain

import "time"

// RunStatus is the lifecycle of an episode production run.
type RunStatus string

const (
	RunStatusPending          RunStatus = "pending"
	RunStatusRunning          RunStatus = "running"
	RunStatusPaused           RunStatus = "paused"
	RunStatusAwaitingApproval RunStatus = "awaiting_approval"
	RunStatusCompleted        RunStatus = "completed"
	RunStatusFailed           RunStatus = "failed"
	RunStatusStopped          RunStatus = "stopped"
)

// ProductionRun is one durable execution producing (or re-producing) an episode.
type ProductionRun struct {
	ID           string                 `json:"id"`
	ProductionID string                 `json:"production_id"`
	ProjectID    string                 `json:"project_id"`
	EpisodeID    string                 `json:"episode_id"`
	Stage        ProductionStage        `json:"stage"`
	Status       RunStatus              `json:"status"`
	BibleVersion int                    `json:"bible_version"`
	WorkflowID   string                 `json:"workflow_id,omitempty"`
	Snapshot     map[string]interface{} `json:"snapshot,omitempty"`
	Jobs         []ProductionJob        `json:"jobs,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	CompletedAt  *time.Time             `json:"completed_at,omitempty"`
}

// CanTransitionTo enforces the run state machine (§37).
func (r *ProductionRun) CanTransitionTo(next RunStatus) bool {
	switch r.Status {
	case RunStatusPending:
		return next == RunStatusRunning || next == RunStatusStopped
	case RunStatusRunning:
		return next == RunStatusPaused || next == RunStatusAwaitingApproval ||
			next == RunStatusCompleted || next == RunStatusFailed || next == RunStatusStopped
	case RunStatusPaused:
		return next == RunStatusRunning || next == RunStatusStopped
	case RunStatusAwaitingApproval:
		return next == RunStatusRunning || next == RunStatusFailed || next == RunStatusStopped
	case RunStatusFailed:
		return next == RunStatusPending // retry re-queues the run
	case RunStatusCompleted, RunStatusStopped:
		return false
	}
	return false
}
