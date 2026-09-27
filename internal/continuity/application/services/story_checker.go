package services

import "context"

type StoryChecker struct{}

func (c *StoryChecker) Check(ctx context.Context, episodeID string) error {
	return nil
}
