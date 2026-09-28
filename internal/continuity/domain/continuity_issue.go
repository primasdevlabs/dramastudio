package domain

import "time"

type Severity string

const (
	SeverityInfo     Severity = "INFO"
	SeverityWarning  Severity = "WARNING"
	SeverityError    Severity = "ERROR"
	SeverityBlocking Severity = "BLOCKING"
)

type IssueStatus string

const (
	IssueOpen     IssueStatus = "open"
	IssueResolved IssueStatus = "resolved"
	IssueWontFix  IssueStatus = "wontfix"
)

// ContinuityIssue is an actionable finding with the state needed for a
// human to resolve it (§38).
type ContinuityIssue struct {
	ID            string      `json:"id"`
	CheckID       string      `json:"check_id,omitempty"`
	ProjectID     string      `json:"project_id"`
	EpisodeID     string      `json:"episode_id,omitempty"`
	SceneID       string      `json:"scene_id,omitempty"`
	Category      string      `json:"category"` // story|character|visual|timeline|wardrobe|production
	Severity      Severity    `json:"severity"`
	Entity        string      `json:"entity"`
	ExpectedState string      `json:"expected_state"`
	ActualState   string      `json:"actual_state"`
	Cause         string      `json:"cause"`
	Evidence      string      `json:"evidence"`
	Resolution    string      `json:"resolution,omitempty"`
	Status        IssueStatus `json:"status"`
	CreatedAt     time.Time   `json:"created_at"`
	ResolvedAt    *time.Time  `json:"resolved_at,omitempty"`
}
