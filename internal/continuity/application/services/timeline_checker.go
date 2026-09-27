package services

import "context"

type TimelineChecker struct{}

func (c *TimelineChecker) Check(ctx context.Context, episodeID string) error {
	return nil
}
