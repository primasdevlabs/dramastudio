package domain

type Character struct {
	ID            CharacterID    `json:"id"`
	Name          string         `json:"name"`
	IsLocked      bool           `json:"is_locked"`
	Appearance    Appearance     `json:"appearance"`
	Personality   Personality    `json:"personality"`
	Wardrobe      Wardrobe       `json:"wardrobe"`
	VoiceProfile  VoiceProfile   `json:"voice_profile"`
	Relationships []Relationship `json:"relationships"`
}

func NewCharacter(id CharacterID, name string) *Character {
	return &Character{
		ID:   id,
		Name: name,
	}
}
