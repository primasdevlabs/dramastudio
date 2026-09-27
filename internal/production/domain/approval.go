package domain

import "time"

type ApprovalDecision string

const (
	DecisionApprove         ApprovalDecision = "APPROVE"
	DecisionReject          ApprovalDecision = "REJECT"
	DecisionRequestRevision ApprovalDecision = "REQUEST_REVISION"
	DecisionRegenerate      ApprovalDecision = "REGENERATE"
	DecisionPause           ApprovalDecision = "PAUSE"
	DecisionStop            ApprovalDecision = "STOP"
	DecisionOverride        ApprovalDecision = "OVERRIDE"
)

type ApprovalRequest struct {
	ID          string           `json:"id"`
	ProjectID   string           `json:"project_id"`
	EpisodeID   string           `json:"episode_id"`
	Stage       string           `json:"stage"` // Bible, Script, Storyboard, Shot, Episode
	TargetID    string           `json:"target_id"`
	Decision    ApprovalDecision `json:"decision"`
	Notes       string           `json:"notes"`
	DecidedBy   string           `json:"decided_by"` // USER, LEAD_DIRECTOR, etc.
	SubmittedAt time.Time        `json:"submitted_at"`
	DecidedAt   *time.Time       `json:"decided_at,omitempty"`
}
