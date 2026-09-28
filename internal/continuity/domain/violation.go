package domain

// Violation is a single rule breach found by a checker, before it is
// persisted as a ContinuityIssue.
type Violation struct {
	RuleName    string   `json:"rule_name"`
	Description string   `json:"description"`
	Severity    Severity `json:"severity"`
	Entity      string   `json:"entity"`
	Evidence    string   `json:"evidence"`
}
