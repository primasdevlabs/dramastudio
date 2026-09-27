package events

import "time"

type Event interface {
	EventType() string
	OccurredAt() time.Time
}

type BaseEvent struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

func (e BaseEvent) EventType() string    { return e.Type }
func (e BaseEvent) OccurredAt() time.Time { return e.Timestamp }

type SeriesCreated struct{ BaseEvent }
type SeasonCreated struct{ BaseEvent }
type ArcApproved struct{ BaseEvent }
type EpisodeCreated struct{ BaseEvent }
type EpisodeApproved struct{ BaseEvent }
type EpisodeCompleted struct{ BaseEvent }
type CharacterCreated struct{ BaseEvent }
type CharacterUpdated struct{ BaseEvent }
type CharacterLocked struct{ BaseEvent }
type StoryFactEstablished struct{ BaseEvent }
type StoryFactChanged struct{ BaseEvent }
type SceneCreated struct{ BaseEvent }
type SceneApproved struct{ BaseEvent }
type AssetGenerated struct{ BaseEvent }
type AssetApproved struct{ BaseEvent }
type ProductionJobCreated struct{ BaseEvent }
type ProductionJobCompleted struct{ BaseEvent }
type ProductionJobFailed struct{ BaseEvent }
type ContinuityIssueDetected struct{ BaseEvent }
type ContinuityIssueResolved struct{ BaseEvent }
type AgentTaskCreated struct{ BaseEvent }
type AgentTaskCompleted struct{ BaseEvent }
type AgentDecisionMade struct{ BaseEvent }
type ApprovalRequested struct{ BaseEvent }
type ApprovalGranted struct{ BaseEvent }
type ApprovalRejected struct{ BaseEvent }
