package http

import (
	"time"

	"dramastudio/internal/characters/domain"
)

// CharacterResponse is the domain-safe shape returned to clients; it never
// exposes internal persistence details (§56).
type CharacterResponse struct {
	ID            string                `json:"id"`
	ProjectID     string                `json:"project_id"`
	Name          string                `json:"name"`
	Role          string                `json:"role"`
	Bio           string                `json:"bio"`
	Version       int                   `json:"version"`
	IsLocked      bool                  `json:"is_locked"`
	Appearance    domain.Appearance     `json:"appearance"`
	Personality   domain.Personality    `json:"personality"`
	VoiceProfile  domain.VoiceProfile   `json:"voice_profile"`
	Wardrobe      domain.Wardrobe       `json:"wardrobe"`
	Relationships []domain.Relationship `json:"relationships"`
}

func toCharacterResponse(c *domain.Character) CharacterResponse {
	return CharacterResponse{
		ID:            string(c.ID),
		ProjectID:     c.ProjectID,
		Name:          c.Name,
		Role:          c.Role,
		Bio:           c.Bio,
		Version:       c.Version,
		IsLocked:      c.IsLocked,
		Appearance:    c.Appearance,
		Personality:   c.Personality,
		VoiceProfile:  c.VoiceProfile,
		Wardrobe:      c.Wardrobe,
		Relationships: c.Relationships,
	}
}

type WardrobeAssignmentResponse struct {
	ID          string                `json:"id"`
	CharacterID string                `json:"character_id"`
	EpisodeID   string                `json:"episode_id"`
	SceneID     string                `json:"scene_id"`
	Items       []domain.WardrobeItem `json:"items"`
	ChangeEvent string                `json:"change_event"`
	CreatedAt   time.Time             `json:"created_at,omitempty"`
}
