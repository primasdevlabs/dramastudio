package domain

import "time"

type FactStatus string

const (
	FactStatusCanonical FactStatus = "canonical"
	FactStatusDisputed  FactStatus = "disputed"
	FactStatusRetconned FactStatus = "retconned"
)

// StoryFact is a canonical piece of story truth (§12).
type StoryFact struct {
	ID                string     `json:"id"`
	ProjectID         string     `json:"project_id"`
	EntityID          string     `json:"entity_id,omitempty"` // subject entity reference
	Subject           string     `json:"subject"`
	Predicate         string     `json:"predicate"`
	Object            string     `json:"object"`
	Type              string     `json:"type,omitempty"`
	IntroducedEpisode string     `json:"introduced_episode,omitempty"`
	EffectiveFrom     string     `json:"effective_from,omitempty"`
	EffectiveUntil    string     `json:"effective_until,omitempty"`
	Source            string     `json:"source,omitempty"`
	Confidence        float64    `json:"confidence"`
	Status            FactStatus `json:"status"`
	Version           int        `json:"version"`
	CreatedAt         time.Time  `json:"created_at"`
}

// IsActive reports whether the fact currently holds (not retconned/expired).
func (f *StoryFact) IsActive() bool {
	return f.Status == FactStatusCanonical
}
