package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"dramastudio/internal/platform/events"
	"dramastudio/internal/platform/storage"
	"dramastudio/internal/postproduction/domain"
)

// Renderer executes a timeline render to a local output path.
type Renderer interface {
	RenderTimeline(ctx context.Context, timeline *domain.Timeline, outputPath string) (string, error)
}

type PostproductionService struct {
	repo     domain.PostproductionRepository
	renderer Renderer
	store    storage.ObjectStorage
	keys     storage.Keys
	workDir  string
	events   *events.Bus // may be nil; set via SetEvents
}

func NewPostproductionService(repo domain.PostproductionRepository, renderer Renderer, store storage.ObjectStorage, workDir string) *PostproductionService {
	return &PostproductionService{repo: repo, renderer: renderer, store: store, keys: storage.NewKeys(), workDir: workDir}
}

// SetEvents injects the domain event bus (§48). Nil-safe emitter.
func (s *PostproductionService) SetEvents(b *events.Bus) {
	s.events = b
}

// SaveTimeline stores a new immutable timeline version for the episode.
func (s *PostproductionService) SaveTimeline(ctx context.Context, projectID, episodeID string, videoTracks, audioTracks []domain.TrackItem, subtitles []domain.Subtitle) (*domain.Timeline, error) {
	version := 1
	if latest, err := s.repo.FindTimelineByEpisode(ctx, episodeID); err == nil {
		version = latest.Version + 1
	}
	tl := &domain.Timeline{
		ID:          "tl_" + uuid.NewString(),
		ProjectID:   projectID,
		EpisodeID:   episodeID,
		Version:     version,
		Status:      domain.TimelineDraft,
		VideoTracks: videoTracks,
		AudioTracks: audioTracks,
		Subtitles:   subtitles,
		CreatedAt:   time.Now().UTC(),
	}
	if err := s.repo.SaveTimeline(ctx, tl); err != nil {
		return nil, err
	}
	return tl, nil
}

func (s *PostproductionService) GetTimeline(ctx context.Context, episodeID string) (*domain.Timeline, error) {
	return s.repo.FindTimelineByEpisode(ctx, episodeID)
}

func (s *PostproductionService) ListTimelineVersions(ctx context.Context, episodeID string) ([]*domain.Timeline, error) {
	return s.repo.ListTimelineVersions(ctx, episodeID)
}

// ApproveTimeline locks a timeline for rendering (human control).
func (s *PostproductionService) ApproveTimeline(ctx context.Context, episodeID string) (*domain.Timeline, error) {
	tl, err := s.repo.FindTimelineByEpisode(ctx, episodeID)
	if err != nil {
		return nil, err
	}
	tl.Status = domain.TimelineApproved
	if err := s.repo.SaveTimeline(ctx, tl); err != nil {
		return nil, err
	}
	return tl, nil
}

// QueueRender creates a queued render task for the latest timeline.
func (s *PostproductionService) QueueRender(ctx context.Context, projectID, episodeID, format, resolution string) (*domain.RenderTask, error) {
	tl, err := s.repo.FindTimelineByEpisode(ctx, episodeID)
	if err != nil {
		return nil, err
	}
	if tl.Status == domain.TimelineDraft {
		return nil, fmt.Errorf("timeline for episode %s is still a draft; approve before rendering", episodeID)
	}
	task := &domain.RenderTask{
		ID:         "rend_" + uuid.NewString(),
		TimelineID: tl.ID,
		ProjectID:  projectID,
		EpisodeID:  episodeID,
		Format:     format,
		Resolution: resolution,
		Status:     domain.RenderQueued,
		CreatedAt:  time.Now().UTC(),
	}
	if err := s.repo.SaveRender(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

// ExecuteRender runs the render synchronously: ffmpeg → object storage →
// render row + timeline status updates. Callers needing async execution
// should invoke this from a Temporal activity.
func (s *PostproductionService) ExecuteRender(ctx context.Context, renderID string) (*domain.RenderTask, error) {
	task, err := s.repo.FindRenderByID(ctx, renderID)
	if err != nil {
		return nil, err
	}
	tl, err := s.repo.FindTimelineByEpisode(ctx, task.EpisodeID)
	if err != nil {
		return nil, err
	}
	task.Status = domain.RenderRendering
	_ = s.repo.SaveRender(ctx, task)

	fail := func(err error) (*domain.RenderTask, error) {
		now := time.Now().UTC()
		task.Status = domain.RenderFailed
		task.Error = err.Error()
		task.FinishedAt = &now
		_ = s.repo.SaveRender(ctx, task)
		s.events.Emit(ctx, events.RenderFailed, task.ProjectID, task.ID,
			fmt.Sprintf("Render failed: %s", err))
		return nil, err
	}

	format := task.Format
	if format == "" {
		format = "mp4"
	}
	workDir := s.workDir
	if workDir == "" {
		workDir = os.TempDir()
	}
	localPath := filepath.Join(workDir, task.ID+"."+format)

	outPath, err := s.renderer.RenderTimeline(ctx, tl, localPath)
	if err != nil {
		return fail(err)
	}

	// Upload to object storage when configured; otherwise keep the local path.
	objectKey := s.keys.Render(task.ProjectID, task.EpisodeID, task.ID)
	outputURL := outPath
	if s.store != nil {
		f, err := os.Open(outPath)
		if err != nil {
			return fail(err)
		}
		info, _ := f.Stat()
		if err := s.store.Upload(ctx, objectKey, f, info.Size(), "video/"+format); err != nil {
			f.Close()
			return fail(err)
		}
		f.Close()
		if url, err := s.store.SignedURL(ctx, objectKey, 24*time.Hour); err == nil {
			outputURL = url
		}
	}

	now := time.Now().UTC()
	task.Status = domain.RenderCompleted
	task.ObjectKey = objectKey
	task.OutputURL = outputURL
	task.FinishedAt = &now
	if err := s.repo.SaveRender(ctx, task); err != nil {
		return nil, err
	}
	tl.Status = domain.TimelineRendered
	_ = s.repo.SaveTimeline(ctx, tl)
	s.events.Emit(ctx, events.RenderCompleted, task.ProjectID, task.ID,
		fmt.Sprintf("Render completed for episode %s", task.EpisodeID))
	return task, nil
}

func (s *PostproductionService) GetRender(ctx context.Context, id string) (*domain.RenderTask, error) {
	return s.repo.FindRenderByID(ctx, id)
}

func (s *PostproductionService) ListRenders(ctx context.Context, episodeID string) ([]*domain.RenderTask, error) {
	return s.repo.ListRendersByEpisode(ctx, episodeID)
}
