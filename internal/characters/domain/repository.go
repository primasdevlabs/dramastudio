package domain

import "context"

type CharacterRepository interface {
	FindByID(ctx context.Context, id CharacterID) (*Character, error)
	List(ctx context.Context) ([]*Character, error)
	Save(ctx context.Context, character *Character) error
}
