// Package contracts holds the vendor-neutral workflow types: the Engine
// port, signal/query names, task-queue names, and the produce-episode
// payload. It is a leaf package so both bounded contexts and engine
// adapters can depend on it without import cycles.
package contracts

import "context"

// Signal and query channels for human control over a running production
// (§39, §83).
const (
	SignalApproval = "approval" // payload: string decision APPROVE|REJECT|REQUEST_REVISION|REGENERATE|OVERRIDE
	SignalPause    = "pause"
	SignalResume   = "resume"
	SignalStop     = "stop"
	QueryStatus    = "status" // returns WorkflowStatus snapshot
	QueryStage     = "stage"
)

// Task queue isolation per Backend.md §83: workloads on separate queues so
// heavy video jobs cannot starve story or QA work.
const (
	QueueStory      = "story-tasks"
	QueueVisual     = "visual-tasks"
	QueueVideo      = "video-tasks"
	QueueAudio      = "audio-tasks"
	QueueMedia      = "media-tasks"
	QueueQA         = "qa-tasks"
	QueuePublishing = "publishing-tasks"
	// QueueCore hosts orchestration workflows (episode production).
	QueueCore = "core-tasks"
)

// AllQueues enumerates every task queue a worker can serve.
var AllQueues = []string{
	QueueCore, QueueStory, QueueVisual, QueueVideo, QueueAudio, QueueMedia, QueueQA, QueuePublishing,
}

// Engine is the workflow-orchestration port (§83). Application code starts,
// signals, queries, and cancels runs through this interface and never sees
// Temporal, Hatchet, or any other vendor API. Concrete engines live under
// platform/workflow/<engine>/ and are selected by WORKFLOW_ENGINE.
type Engine interface {
	// Start launches an episode production run. The returned workflow ID
	// is stored on the ProductionRun for later signals/queries.
	Start(ctx context.Context, in ProduceEpisodeInput) (workflowID string, err error)

	// Signal delivers a control signal (SignalApproval/Pause/Resume/Stop)
	// to a running workflow.
	Signal(ctx context.Context, workflowID, signal string, payload interface{}) error

	// Cancel force-terminates a workflow.
	Cancel(ctx context.Context, workflowID string) error

	// Status returns the engine's view of a run for the studio UI.
	Status(ctx context.Context, workflowID string) (*WorkflowStatus, error)

	// Close releases engine resources.
	Close() error
}

// ProduceEpisodeInput starts one episode production run.
type ProduceEpisodeInput struct {
	ProjectID    string `json:"project_id"`
	EpisodeID    string `json:"episode_id"`
	RunID        string `json:"run_id"`
	BibleVersion int    `json:"bible_version"`
}

// ProduceEpisodeResult is the durable outcome of a completed run.
type ProduceEpisodeResult struct {
	EpisodeID      string   `json:"episode_id"`
	RenderURL      string   `json:"render_url"`
	Approval       string   `json:"approval"`
	Status         string   `json:"status"`
	GeneratedShots []string `json:"generated_shots"`
	Plan           string   `json:"plan"`
}

// WorkflowStatus is the query surface for the studio UI.
type WorkflowStatus struct {
	RunID      string `json:"run_id"`
	Stage      string `json:"stage"`
	Paused     bool   `json:"paused"`
	Stopped    bool   `json:"stopped"`
	ApprovalID string `json:"approval_id,omitempty"`
}
