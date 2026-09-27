package domain

type TaskStatus string

const (
	TaskPending   TaskStatus = "pending"
	TaskExecuting TaskStatus = "executing"
	TaskCompleted TaskStatus = "completed"
	TaskFailed    TaskStatus = "failed"
)
