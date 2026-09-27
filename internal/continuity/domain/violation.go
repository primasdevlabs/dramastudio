package domain

type Violation struct {
	ID          string   `json:"id"`
	RuleName    string   `json:"rule_name"`
	Description string   `json:"description"`
	Severity    Severity `json:"severity"`
}
