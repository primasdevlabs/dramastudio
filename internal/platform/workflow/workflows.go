package workflow

import (
	"context"
	"fmt"
	"log"

	"dramastudio/internal/platform/workflow/activities"
)

type ProduceEpisodeResult struct {
	EpisodeID    string   `json:"episode_id"`
	RenderURL    string   `json:"render_url"`
	Approval     string   `json:"approval"`
	Status       string   `json:"status"`
	GeneratedShots []string `json:"generated_shots"`
}

type ProduceEpisodeWorkflow struct {
	act *activities.EpisodeActivities
}

func NewProduceEpisodeWorkflow(act *activities.EpisodeActivities) *ProduceEpisodeWorkflow {
	if act == nil {
		act = activities.NewEpisodeActivities()
	}
	return &ProduceEpisodeWorkflow{act: act}
}

func (w *ProduceEpisodeWorkflow) Execute(ctx context.Context, projectID, episodeID string) (*ProduceEpisodeResult, error) {
	log.Printf("[Workflow] Starting ProduceEpisodeWorkflow project=%s episode=%s", projectID, episodeID)

	// 1. Episode Plan
	plan, err := w.act.GenerateEpisodePlan(ctx, projectID, episodeID)
	if err != nil {
		return nil, fmt.Errorf("GenerateEpisodePlan failed: %w", err)
	}

	// 2. Validate Continuity
	ok, err := w.act.ValidateContinuity(ctx, projectID, episodeID)
	if err != nil || !ok {
		return nil, fmt.Errorf("ValidateContinuity failed: %w", err)
	}

	// 3. Script & Dialogue
	_, err = w.act.GenerateScript(ctx, episodeID)
	if err != nil {
		return nil, fmt.Errorf("GenerateScript failed: %w", err)
	}

	_, err = w.act.GenerateDialogue(ctx, episodeID)
	if err != nil {
		return nil, fmt.Errorf("GenerateDialogue failed: %w", err)
	}

	// 4. Scene Plans & Storyboards
	scenes, err := w.act.GenerateScenePlans(ctx, episodeID)
	if err != nil || len(scenes) == 0 {
		return nil, fmt.Errorf("GenerateScenePlans failed: %w", err)
	}

	shots, err := w.act.GenerateStoryboard(ctx, scenes[0])
	if err != nil {
		return nil, fmt.Errorf("GenerateStoryboard failed: %w", err)
	}

	// 5. Generate Video Shots (Wan)
	videoURLs, err := w.act.GenerateVideoShots(ctx, shots)
	if err != nil {
		return nil, fmt.Errorf("GenerateVideoShots failed: %w", err)
	}

	// 6. Assemble Episode
	renderURL, err := w.act.AssembleEpisode(ctx, episodeID, videoURLs)
	if err != nil {
		return nil, fmt.Errorf("AssembleEpisode failed: %w", err)
	}

	// 7. Request Human Approval
	approval, err := w.act.RequestHumanApproval(ctx, episodeID)
	if err != nil {
		return nil, fmt.Errorf("RequestHumanApproval failed: %w", err)
	}

	_ = plan
	return &ProduceEpisodeResult{
		EpisodeID:      episodeID,
		RenderURL:      renderURL,
		Approval:       approval,
		Status:         "COMPLETED",
		GeneratedShots: videoURLs,
	}, nil
}
