package domain

type AgentStatus string

const (
	StatusIdle        AgentStatus = "idle"
	StatusWorking     AgentStatus = "working"
	StatusInterrupted AgentStatus = "interrupted"
)
