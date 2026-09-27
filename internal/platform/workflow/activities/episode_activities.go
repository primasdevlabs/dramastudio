package activities

import (
	"context"
	"fmt"
	"log"
)

type EpisodeActivities struct{}

func NewEpisodeActivities() *EpisodeActivities {
	return &EpisodeActivities{}
}

func (a *EpisodeActivities) GenerateEpisodePlan(ctx context.Context, projectID, episodeID string) (string, error) {
	log.Printf("[Activity] GenerateEpisodePlan project=%s episode=%s", projectID, episodeID)
	return fmt.Sprintf("plan_%s", episodeID), nil
}

func (a *EpisodeActivities) ValidateContinuity(ctx context.Context, projectID, episodeID string) (bool, error) {
	log.Printf("[Activity] ValidateContinuity project=%s episode=%s", projectID, episodeID)
	return true, nil
}

func (a *EpisodeActivities) GenerateScript(ctx context.Context, episodeID string) (string, error) {
	log.Printf("[Activity] GenerateScript episode=%s", episodeID)
	return fmt.Sprintf("script_%s", episodeID), nil
}

func (a *EpisodeActivities) GenerateDialogue(ctx context.Context, episodeID string) (string, error) {
	log.Printf("[Activity] GenerateDialogue episode=%s", episodeID)
	return fmt.Sprintf("dialogue_%s", episodeID), nil
}

func (a *EpisodeActivities) GenerateScenePlans(ctx context.Context, episodeID string) ([]string, error) {
	log.Printf("[Activity] GenerateScenePlans episode=%s", episodeID)
	return []string{fmt.Sprintf("scene_01_%s", episodeID)}, nil
}

func (a *EpisodeActivities) GenerateStoryboard(ctx context.Context, sceneID string) ([]string, error) {
	log.Printf("[Activity] GenerateStoryboard scene=%s", sceneID)
	return []string{"shot_001", "shot_002"}, nil
}

func (a *EpisodeActivities) GenerateVideoShots(ctx context.Context, shotIDs []string) ([]string, error) {
	log.Printf("[Activity] GenerateVideoShots shots=%v", shotIDs)
	urls := make([]string, len(shotIDs))
	for i, id := range shotIDs {
		urls[i] = fmt.Sprintf("https://storage.dramastudio.ai/renders/wan_%s.mp4", id)
	}
	return urls, nil
}

func (a *EpisodeActivities) AssembleEpisode(ctx context.Context, episodeID string, videoURLs []string) (string, error) {
	log.Printf("[Activity] AssembleEpisode episode=%s URLs=%v", episodeID, videoURLs)
	return fmt.Sprintf("https://storage.dramastudio.ai/renders/ep_%s_render.mp4", episodeID), nil
}

func (a *EpisodeActivities) RequestHumanApproval(ctx context.Context, episodeID string) (string, error) {
	log.Printf("[Activity] RequestHumanApproval episode=%s", episodeID)
	return "APPROVED", nil
}
