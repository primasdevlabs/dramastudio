package domain

type CanonRule struct {
	ID          string `json:"id"`
	ProjectID   string `json:"project_id"`
	RuleText    string `json:"rule_text"`
	Enforced    bool   `json:"enforced"`
}
