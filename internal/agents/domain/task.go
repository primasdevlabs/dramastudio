package domain

type Task struct {
	ID          string     `json:"id"`
	AgentID     AgentID    `json:"agent_id"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
}
