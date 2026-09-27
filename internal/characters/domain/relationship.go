package domain

type Relationship struct {
	TargetCharacterID CharacterID `json:"target_character_id"`
	RelationshipType  string      `json:"relationship_type"`
	Description       string      `json:"description"`
}
