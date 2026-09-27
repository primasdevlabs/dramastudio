package services

import "context"

type WardrobeChecker struct{}

func (c *WardrobeChecker) Check(ctx context.Context, characterID string) error {
	return nil
}
