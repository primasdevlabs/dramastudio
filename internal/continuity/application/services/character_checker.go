package services

import "context"

type CharacterChecker struct{}

func (c *CharacterChecker) Check(ctx context.Context, characterID string) error {
	return nil
}
