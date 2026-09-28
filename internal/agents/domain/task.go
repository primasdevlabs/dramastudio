package domain

import "time"

type TaskStatus string

const (
	TaskPending   TaskStatus = "pending"
	TaskRunning   TaskStatus = "running"
	TaskSucceeded TaskStatus = "succeeded"
	TaskFailed    TaskStatus = "failed"
	TaskEscalated TaskStatus = "escalated"
	TaskCancelled TaskStatus = "cancelled"
)

// ApprovalPolicy controls whether a task executes autonomously or waits
// for human approval (§44).
type ApprovalPolicy string

const (
	PolicyMonitored  ApprovalPolicy = "monitored"
	PolicyAutonomous ApprovalPolicy = "autonomous"
)

// Task is a bounded unit of work delegated to an agent.
type Task struct {
	ID              string                 `json:"id"`
	ProjectID       string                 `json:"project_id"`
	AgentID         string                 `json:"agent_id"`
	Objective       string                 `json:"objective"`
	Input           map[string]interface{} `json:"input"`
	ExpectedOutput  string                 `json:"expected_output"`
	Constraints     map[string]interface{} `json:"constraints"`
	Budget          *float64               `json:"budget,omitempty"`
	MaxIterations   int                    `json:"max_iterations"`
	Iteration       int                    `json:"iteration"`
	TimeoutSec      int                    `json:"timeout_sec"`
	ApprovalPolicy  ApprovalPolicy         `json:"approval_policy"`
	SuccessCriteria string                 `json:"success_criteria"`
	Status          TaskStatus             `json:"status"`
	Result          map[string]interface{} `json:"result,omitempty"`
	Error           string                 `json:"error,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	CompletedAt     *time.Time             `json:"completed_at,omitempty"`
}

// CanTransitionTo enforces the task lifecycle.
func (t *Task) CanTransitionTo(next TaskStatus) bool {
	switch t.Status {
	case TaskPending:
		return next == TaskRunning || next == TaskCancelled || next == TaskEscalated
	case TaskRunning:
		return next == TaskSucceeded || next == TaskFailed || next == TaskEscalated || next == TaskCancelled
	case TaskFailed:
		return next == TaskPending // retry
	case TaskEscalated:
		return next == TaskPending || next == TaskCancelled
	case TaskSucceeded, TaskCancelled:
		return false
	}
	return false
}
