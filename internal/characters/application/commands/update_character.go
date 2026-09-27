package commands

import "dramastudio/internal/characters/domain"

type UpdateCharacter struct {
	ID   domain.CharacterID
	Name string
}
