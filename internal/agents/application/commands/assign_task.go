package commands

import "dramastudio/internal/agents/domain"

type AssignTask struct {
	AgentID     domain.AgentID
	Description string
}
