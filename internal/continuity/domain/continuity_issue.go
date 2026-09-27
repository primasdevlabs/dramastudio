package domain

type Severity string

const (
	SeverityInfo     Severity = "INFO"
	SeverityWarning  Severity = "WARNING"
	SeverityError    Severity = "ERROR"
	SeverityBlocking Severity = "BLOCKING"
)

type ContinuityIssue struct {
	ID            string   `json:"id"`
	ProjectID     string   `json:"project_id"`
	EpisodeID     string   `json:"episode_id"`
	SceneID       string   `json:"scene_id"`
	Category      string   `json:"category"` // Story, Character, Visual, Timeline, Wardrobe
	Severity      Severity `json:"severity"`
	Entity        string   `json:"entity"`
	ExpectedState string   `json:"expected_state"`
	ActualState   string   `json:"actual_state"`
	Cause         string   `json:"cause"`
	Evidence      string   `json:"evidence"`
	Resolution    string   `json:"resolution"`
	IsResolved    bool     `json:"is_resolved"`
}
