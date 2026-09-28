package activities

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	aiapp "dramastudio/internal/ai/application"
	"dramastudio/internal/ai/domain"
	continuitysvc "dramastudio/internal/continuity/application/services"
	continuitydomain "dramastudio/internal/continuity/domain"
	mediasvc "dramastudio/internal/media/application/services"
	mediadomain "dramastudio/internal/media/domain"
	postdomain "dramastudio/internal/postproduction/domain"
	proddomain "dramastudio/internal/production/domain"
)

// EpisodeActivities implements the activity set orchestrated by
// ProduceEpisodeWorkflow (§41).
type EpisodeActivities struct {
	deps Deps
}

func NewEpisodeActivities(deps Deps) *EpisodeActivities {
	return &EpisodeActivities{deps: deps}
}

// runAI executes a text-generation capability via the ai context and
// returns the model's text output.
func (a *EpisodeActivities) runAI(ctx context.Context, projectID, scope string, capability domain.AICapability, prompt string, input interface{}) (string, error) {
	if a.deps.AI == nil {
		return "", fmt.Errorf("ai executor not configured")
	}
	raw, _ := json.Marshal(input)
	job, err := a.deps.AI.Handle(ctx, aiapp.ExecuteGenerationCommand{
		ProjectID:  projectID,
		Capability: capability,
		Scope:      scope,
		Input:      raw,
		Prompt:     prompt,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", capability, err)
	}
	var out struct {
		Text string `json:"text"`
	}
	_ = json.Unmarshal(job.Output, &out)
	return out.Text, nil
}

func (a *EpisodeActivities) GenerateEpisodePlan(ctx context.Context, projectID, episodeID string) (string, error) {
	return a.runAI(ctx, projectID, "episode:"+episodeID, domain.CapabilityEpisode,
		fmt.Sprintf("Produce an episode plan for episode %s", episodeID),
		map[string]string{"episode_id": episodeID})
}

// ValidateContinuity runs the timeline check suite and reports the number
// of blocking issues found (workflow fails if > 0).
func (a *EpisodeActivities) ValidateContinuity(ctx context.Context, projectID, episodeID string) (int, error) {
	if a.deps.Continuity == nil {
		return 0, fmt.Errorf("continuity service not configured")
	}
	_, issues, err := a.deps.Continuity.RunCheck(ctx, continuitysvc.CheckInput{
		ProjectID: projectID,
		EpisodeID: episodeID,
		CheckType: continuitydomain.CheckTimeline,
	})
	if err != nil {
		return 0, err
	}
	blocking := 0
	for _, i := range issues {
		if i.Severity == continuitydomain.SeverityBlocking {
			blocking++
		}
	}
	return blocking, nil
}

func (a *EpisodeActivities) GenerateScript(ctx context.Context, projectID, episodeID string) (string, error) {
	return a.runAI(ctx, projectID, "episode:"+episodeID, domain.CapabilityScript,
		fmt.Sprintf("Write the script for episode %s", episodeID),
		map[string]string{"episode_id": episodeID})
}

func (a *EpisodeActivities) GenerateDialogue(ctx context.Context, projectID, episodeID string) (string, error) {
	return a.runAI(ctx, projectID, "episode:"+episodeID, domain.CapabilityDialogue,
		fmt.Sprintf("Write dialogue for episode %s", episodeID),
		map[string]string{"episode_id": episodeID})
}

// GenerateScenePlans returns the episode's scene IDs; the story module is
// authoritative, so this reads rather than inventing scene structure.
func (a *EpisodeActivities) GenerateScenePlans(ctx context.Context, episodeID string) ([]string, error) {
	if a.deps.Story == nil {
		return nil, fmt.Errorf("story service not configured")
	}
	scenes, err := a.deps.Story.ListScenes(ctx, episodeID)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(scenes))
	for _, s := range scenes {
		ids = append(ids, s.ID)
	}
	return ids, nil
}

// GenerateStoryboard produces the storyboard spec per scene and records a
// planned shot per scene.
func (a *EpisodeActivities) GenerateStoryboard(ctx context.Context, projectID, episodeID string, sceneIDs []string) ([]string, error) {
	if a.deps.Production == nil {
		return nil, fmt.Errorf("production service not configured")
	}
	shotIDs := make([]string, 0, len(sceneIDs))
	for i, sceneID := range sceneIDs {
		desc, err := a.runAI(ctx, projectID, "episode:"+episodeID, domain.CapabilityStoryboard,
			fmt.Sprintf("Storyboard scene %s of episode %s", sceneID, episodeID),
			map[string]string{"episode_id": episodeID, "scene_id": sceneID})
		if err != nil {
			return nil, err
		}
		shot, err := a.deps.Production.CreateShot(ctx, &proddomain.Shot{
			ProjectID:   projectID,
			EpisodeID:   episodeID,
			SceneID:     sceneID,
			Seq:         i + 1,
			Description: desc,
		})
		if err != nil {
			return nil, err
		}
		shotIDs = append(shotIDs, shot.ID)
	}
	return shotIDs, nil
}

// GenerateVideoShots submits video generation for each shot and returns
// the asset IDs. Async provider completion arrives via webhook.
func (a *EpisodeActivities) GenerateVideoShots(ctx context.Context, projectID, episodeID string, shotIDs []string) ([]string, error) {
	if a.deps.Media == nil {
		return nil, fmt.Errorf("media service not configured")
	}
	assetIDs := make([]string, 0, len(shotIDs))
	for _, shotID := range shotIDs {
		shot, err := a.deps.Production.GetShot(ctx, shotID)
		if err != nil {
			return nil, err
		}
		_, asset, err := a.deps.Media.Generate(ctx, mediasvc.GenerateRequest{
			ProjectID:   projectID,
			Capability:  "video_generation",
			MediaType:   mediadomain.MediaTypeVideo,
			EpisodeID:   episodeID,
			SceneID:     shot.SceneID,
			ShotID:      shotID,
			Prompt:      shot.Description,
			CallbackURL: a.deps.CallbackURL,
		})
		if err != nil {
			return nil, err
		}
		assetIDs = append(assetIDs, asset.ID)
	}
	return assetIDs, nil
}

// AssembleEpisode builds the render timeline from approved shot assets and
// executes the render synchronously (a heavyweight deploy would split this
// into queued render workers).
func (a *EpisodeActivities) AssembleEpisode(ctx context.Context, projectID, episodeID string, shotIDs []string) (string, error) {
	if a.deps.Postprod == nil {
		return "", fmt.Errorf("postproduction service not configured")
	}
	start := 0.0
	tracks := make([]postdomain.TrackItem, 0, len(shotIDs))
	for _, shotID := range shotIDs {
		shot, err := a.deps.Production.GetShot(ctx, shotID)
		if err != nil {
			return "", err
		}
		if shot.ApprovedAsset == "" {
			continue // only approved/generated shots land on the timeline
		}
		tracks = append(tracks, postdomain.TrackItem{
			ID:        "ti_" + uuid.NewString(),
			ShotID:    shotID,
			StartTime: start,
			Duration:  shot.DurationSec,
			AssetURL:  shot.ApprovedAsset,
		})
		start += shot.DurationSec
	}
	if len(tracks) == 0 {
		return "", fmt.Errorf("no approved shot assets to assemble for episode %s", episodeID)
	}
	if _, err := a.deps.Postprod.SaveTimeline(ctx, projectID, episodeID, tracks, nil, nil); err != nil {
		return "", err
	}
	if _, err := a.deps.Postprod.ApproveTimeline(ctx, episodeID); err != nil {
		return "", err
	}
	task, err := a.deps.Postprod.QueueRender(ctx, projectID, episodeID, "mp4", "1080x1920")
	if err != nil {
		return "", err
	}
	done, err := a.deps.Postprod.ExecuteRender(ctx, task.ID)
	if err != nil {
		return "", err
	}
	return done.OutputURL, nil
}

// RequestHumanApproval records an approval gate; the workflow itself waits
// on the Temporal signal until a human decides.
func (a *EpisodeActivities) RequestHumanApproval(ctx context.Context, projectID, episodeID, stage, targetID string) (string, error) {
	if a.deps.Production == nil {
		return "", fmt.Errorf("production service not configured")
	}
	req, err := a.deps.Production.RequestApproval(ctx, projectID, episodeID, stage, targetID)
	if err != nil {
		return "", err
	}
	return req.ID, nil
}

// MarkShotGenerating / MarkShotGenerated update shot status around async
// generation so operators see progress in real time.
func (a *EpisodeActivities) MarkShotGenerating(ctx context.Context, shotID string) error {
	_, err := a.deps.Production.SetShotStatus(ctx, shotID, proddomain.ShotStatusGenerating)
	return err
}

func (a *EpisodeActivities) MarkShotGenerated(ctx context.Context, shotID, assetURL string) error {
	shot, err := a.deps.Production.GetShot(ctx, shotID)
	if err != nil {
		return err
	}
	shot.ApprovedAsset = assetURL
	shot.Status = proddomain.ShotStatusGenerated
	_, err = a.deps.Production.UpdateShot(ctx, shot)
	return err
}
