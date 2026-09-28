package events

import "time"

// Versioned domain event vocabulary (§48). These names are emitted on the
// bus, persisted to platform.domain_events, and streamed to clients over
// SSE — treat them as a public contract.
const (
	ProjectCreated   = "ProjectCreated"
	SeriesCreated    = "SeriesCreated"
	SeasonCreated    = "SeasonCreated"
	ArcCreated       = "ArcCreated"
	EpisodeCreated   = "EpisodeCreated"
	EpisodeApproved  = "EpisodeApproved"
	EpisodeCompleted = "EpisodeCompleted"

	CharacterCreated = "CharacterCreated"
	CharacterUpdated = "CharacterUpdated"
	CharacterLocked  = "CharacterLocked"

	StoryFactEstablished = "StoryFactEstablished"
	StoryFactChanged     = "StoryFactChanged"

	SceneCreated  = "SceneCreated"
	SceneApproved = "SceneApproved"
	LocationAdded = "LocationAdded"

	GenerationStarted   = "GenerationStarted"
	GenerationCompleted = "GenerationCompleted"
	AssetGenerated      = "AssetGenerated"
	AssetApproved       = "AssetApproved"
	AssetRejected       = "AssetRejected"

	ProductionJobCreated   = "ProductionJobCreated"
	ProductionJobCompleted = "ProductionJobCompleted"
	ProductionJobFailed    = "ProductionJobFailed"
	ProductionRunStarted   = "ProductionRunStarted"
	ProductionRunPaused    = "ProductionRunPaused"
	ProductionRunResumed   = "ProductionRunResumed"
	ProductionRunStopped   = "ProductionRunStopped"

	ContinuityIssueDetected = "ContinuityIssueDetected"
	ContinuityIssueResolved = "ContinuityIssueResolved"

	AgentTaskCreated   = "AgentTaskCreated"
	AgentTaskCompleted = "AgentTaskCompleted"
	AgentTaskFailed    = "AgentTaskFailed"
	AgentDecisionMade  = "AgentDecisionMade"

	ApprovalRequested = "ApprovalRequested"
	ApprovalGranted   = "ApprovalGranted"
	ApprovalRejected  = "ApprovalRejected"

	RenderCompleted      = "RenderCompleted"
	RenderFailed         = "RenderFailed"
	PublicationPublished = "PublicationPublished"
	PublicationFailed    = "PublicationFailed"

	BibleUpdated    = "SeriesBibleUpdated"
	PlotThreadAdded = "PlotThreadCreated"
	ProjectUpdated  = "ProjectUpdated"
	ApprovalDecided = "ApprovalDecided"
)

// DomainEvent is the concrete event carried by the bus. Its fields map
// directly onto the platform.domain_events columns so the durable log and
// the realtime stream share one shape (§48, §58).
type DomainEvent struct {
	Type        string                 `json:"type"`
	Version     int                    `json:"version"`
	Aggregate   string                 `json:"aggregate,omitempty"`
	AggregateID string                 `json:"aggregate_id,omitempty"`
	ProjectID   string                 `json:"project_id,omitempty"`
	Message     string                 `json:"message,omitempty"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
	OccurredAt  time.Time              `json:"occurred_at"`
}

// New builds a v1 domain event stamped with the current time.
func New(eventType, projectID, aggregateID, message string) DomainEvent {
	return DomainEvent{
		Type:        eventType,
		Version:     1,
		AggregateID: aggregateID,
		ProjectID:   projectID,
		Message:     message,
		OccurredAt:  time.Now().UTC(),
	}
}

// EventType exposes the event name for log/subscription routing.
func (e DomainEvent) EventType() string { return e.Type }
