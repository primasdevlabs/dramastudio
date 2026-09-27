package domain

type CanonRule struct {
	ID          string `json:"id"`
	RuleText    string `json:"rule_text"`
	Enforced    bool   `json:"enforced"`
}
