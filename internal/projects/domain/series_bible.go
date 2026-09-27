package domain

import "time"

type SeriesBible struct {
	ID               string    `json:"id"`
	ProjectID        ProjectID `json:"project_id"`
	Version          int       `json:"version"`
	Premise          string    `json:"premise"`
	Genre            string    `json:"genre"`
	Themes           []string  `json:"themes"`
	Tone             string    `json:"tone"`
	WorldRules       []string  `json:"world_rules"`
	NarrativeRules   []string  `json:"narrative_rules"`
	VisualDirection  string    `json:"visual_direction"`
	DialogueStyle    string    `json:"dialogue_style"`
	StoryConstraints []string  `json:"story_constraints"`
	CreatedAt        time.Time `json:"created_at"`
}

func NewSeriesBible(id string, projectID ProjectID, version int, premise string) *SeriesBible {
	return &SeriesBible{
		ID:               id,
		ProjectID:        projectID,
		Version:          version,
		Premise:          premise,
		Themes:           []string{},
		WorldRules:       []string{},
		NarrativeRules:   []string{},
		StoryConstraints: []string{},
		CreatedAt:        time.Now().UTC(),
	}
}
