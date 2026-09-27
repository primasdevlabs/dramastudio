package domain

type AgentRole string

const (
	RoleLeadDirector AgentRole = "lead_director"
	RoleStoryAgent   AgentRole = "story_agent"
	RoleVisualAgent  AgentRole = "visual_agent"
	RoleQAAgent      AgentRole = "qa_agent"
)
