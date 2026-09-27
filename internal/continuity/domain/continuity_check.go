package domain

type ContinuityCheck struct {
	ID        string    `json:"id"`
	CheckType CheckType `json:"check_type"`
	TargetID  string    `json:"target_id"`
}
