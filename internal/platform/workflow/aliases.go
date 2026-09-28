// Aliases keep the historical "workflow.X" surface working while the
// vendor-neutral types live in the leaf contracts package (which bounded
// contexts may import without cycles).
package workflow

import "dramastudio/internal/platform/workflow/contracts"

type (
	ProduceEpisodeInput  = contracts.ProduceEpisodeInput
	ProduceEpisodeResult = contracts.ProduceEpisodeResult
	WorkflowStatus       = contracts.WorkflowStatus
	Engine               = contracts.Engine
)

const (
	SignalApproval = contracts.SignalApproval
	SignalPause    = contracts.SignalPause
	SignalResume   = contracts.SignalResume
	SignalStop     = contracts.SignalStop
	QueryStatus    = contracts.QueryStatus
	QueryStage     = contracts.QueryStage

	QueueStory      = contracts.QueueStory
	QueueVisual     = contracts.QueueVisual
	QueueVideo      = contracts.QueueVideo
	QueueAudio      = contracts.QueueAudio
	QueueMedia      = contracts.QueueMedia
	QueueQA         = contracts.QueueQA
	QueuePublishing = contracts.QueuePublishing
	QueueCore       = contracts.QueueCore
)

var AllQueues = contracts.AllQueues
