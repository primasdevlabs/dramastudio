package domain

type Decision struct {
	ID         string `json:"id"`
	AgentID    AgentID `json:"agent_id"`
	Rationale  string `json:"rationale"`
	IsApproved bool   `json:"is_approved"`
}
