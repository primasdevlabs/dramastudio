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

// runTransitions encodes the run state machine (§37).
var runTransitions = map[RunStatus][]RunStatus{
	RunStatusPending:          {RunStatusRunning, RunStatusStopped},
	RunStatusRunning:          {RunStatusPaused, RunStatusAwaitingApproval, RunStatusCompleted, RunStatusFailed, RunStatusStopped},
	RunStatusPaused:           {RunStatusRunning, RunStatusStopped},
	RunStatusAwaitingApproval: {RunStatusRunning, RunStatusFailed, RunStatusStopped},
	RunStatusFailed:           {RunStatusPending}, // retry re-queues the run
	RunStatusCompleted:        {},
	RunStatusStopped:          {},
}

// CanTransitionTo reports whether the run may move to next.
func (r *ProductionRun) CanTransitionTo(next RunStatus) bool {
	for _, allowed := range runTransitions[r.Status] {
		if allowed == next {
			return true
		}
	}
	return false
}
