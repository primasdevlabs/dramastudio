package domain

type Character struct {
	ID            CharacterID    `json:"id"`
	ProjectID     string         `json:"project_id"`
	Version       int            `json:"version"`
	Name          string         `json:"name"`
	Role          string         `json:"role"`
	Bio           string         `json:"bio"`
	IsLocked      bool           `json:"is_locked"`
	Appearance    Appearance     `json:"appearance"`
	Personality   Personality    `json:"personality"`
	Wardrobe      Wardrobe       `json:"wardrobe"`
	VoiceProfile  VoiceProfile   `json:"voice_profile"`
	Relationships []Relationship `json:"relationships"`
}

func NewCharacter(id CharacterID, projectID, name, role string) *Character {
	return &Character{
		ID:            id,
		ProjectID:     projectID,
		Version:       1,
		Name:          name,
		Role:          role,
		Relationships: []Relationship{},
	}
}
