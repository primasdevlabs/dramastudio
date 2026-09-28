package domain

import "time"

type CheckStatus string

const (
	CheckRunning  CheckStatus = "running"
	CheckFinished CheckStatus = "finished"
	CheckFailed   CheckStatus = "failed"
)

// ContinuityCheck is one execution of a check suite against a target
// (episode/scene/asset); produced issues link back via CheckID.
type ContinuityCheck struct {
	ID         string      `json:"id"`
	ProjectID  string      `json:"project_id"`
	EpisodeID  string      `json:"episode_id,omitempty"`
	CheckType  CheckType   `json:"check_type"`
	TargetID   string      `json:"target_id,omitempty"`
	Status     CheckStatus `json:"status"`
	IssueCount int         `json:"issue_count"`
	CreatedAt  time.Time   `json:"created_at"`
	FinishedAt *time.Time  `json:"finished_at,omitempty"`
}
