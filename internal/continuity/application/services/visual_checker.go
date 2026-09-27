package services

import "context"

type VisualChecker struct{}

func (c *VisualChecker) Check(ctx context.Context, assetID string) error {
	return nil
}
