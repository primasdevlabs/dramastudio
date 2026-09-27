package http

import (
	"dramastudio/internal/characters/domain"
	platformhttp "dramastudio/internal/platform/http"
)

type createCharacterReq struct {
	Name string `json:"name"`
	Role string `json:"role"`
	Bio  string `json:"bio"`
}

func (r *createCharacterReq) Validate() error {
	if r.Name == "" {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "name", Message: "required"}}}
	}
	return nil
}

type updateCharacterReq struct {
	Name         *string              `json:"name"`
	Role         *string              `json:"role"`
	Bio          *string              `json:"bio"`
	Appearance   *domain.Appearance   `json:"appearance"`
	Personality  *domain.Personality  `json:"personality"`
	VoiceProfile *domain.VoiceProfile `json:"voice_profile"`
	Wardrobe     *domain.Wardrobe     `json:"wardrobe"`
}

func (r *updateCharacterReq) Validate() error { return nil }

type relationshipReq struct {
	TargetCharacterID string `json:"target_character_id"`
	RelationshipType  string `json:"relationship_type"`
	Description       string `json:"description"`
}

func (r *relationshipReq) Validate() error {
	ve := &platformhttp.ValidationError{}
	if r.TargetCharacterID == "" {
		ve.Add("target_character_id", "required")
	}
	if r.RelationshipType == "" {
		ve.Add("relationship_type", "required")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

type wardrobeReq struct {
	EpisodeID   string                `json:"episode_id"`
	SceneID     string                `json:"scene_id"`
	Items       []domain.WardrobeItem `json:"items"`
	ChangeEvent string                `json:"change_event"`
}

func (r *wardrobeReq) Validate() error {
	ve := &platformhttp.ValidationError{}
	if r.EpisodeID == "" {
		ve.Add("episode_id", "required")
	}
	if len(r.Items) == 0 {
		ve.Add("items", "required")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}
