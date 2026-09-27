package domain

import "context"

// CharacterVersion is an immutable snapshot used by production (§15).
type CharacterVersion struct {
	ID           string       `json:"id"`
	CharacterID  CharacterID  `json:"character_id"`
	Version      int          `json:"version"`
	Appearance   Appearance   `json:"appearance"`
	Personality  Personality  `json:"personality"`
	VoiceProfile VoiceProfile `json:"voice_profile"`
	Wardrobe     Wardrobe     `json:"wardrobe"`
}

// WardrobeAssignment ties wardrobe to a production scope for continuity (§17).
type WardrobeAssignment struct {
	ID          string         `json:"id"`
	CharacterID CharacterID    `json:"character_id"`
	EpisodeID   string         `json:"episode_id"`
	SceneID     string         `json:"scene_id"`
	Items       []WardrobeItem `json:"items"`
	ChangeEvent string         `json:"change_event"` // why the wardrobe changed
}

type CharacterRepository interface {
	Save(ctx context.Context, character *Character) error
	FindByID(ctx context.Context, id CharacterID) (*Character, error)
	List(ctx context.Context, projectID string) ([]*Character, error)
	SetLocked(ctx context.Context, id CharacterID, locked bool) error

	SaveVersion(ctx context.Context, v *CharacterVersion) error
	FindVersion(ctx context.Context, characterID CharacterID, version int) (*CharacterVersion, error)
	ListVersions(ctx context.Context, characterID CharacterID) ([]*CharacterVersion, error)

	SaveRelationship(ctx context.Context, characterID CharacterID, rel *Relationship) error

	SaveWardrobeAssignment(ctx context.Context, wa *WardrobeAssignment) error
	ListWardrobeAssignments(ctx context.Context, characterID CharacterID, episodeID string) ([]*WardrobeAssignment, error)
}
