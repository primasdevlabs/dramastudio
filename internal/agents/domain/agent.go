package domain

type Agent struct {
	ID     AgentID     `json:"id"`
	Role   AgentRole   `json:"role"`
	Status AgentStatus `json:"status"`
	Skills []Skill     `json:"skills"`
}
