package workflow

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

// AllQueues enumerates every task queue the worker can serve.
var AllQueues = []string{
	QueueCore, QueueStory, QueueVisual, QueueVideo, QueueAudio, QueueMedia, QueueQA, QueuePublishing,
}
