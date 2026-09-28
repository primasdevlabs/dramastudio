package persistence

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"

	"dramastudio/internal/platform/database/postgres"
	"dramastudio/internal/postproduction/domain"
)

type PostgresPostproductionRepository struct {
	q postgres.Querier
}

func NewPostgresPostproductionRepository(q postgres.Querier) *PostgresPostproductionRepository {
	return &PostgresPostproductionRepository{q: q}
}

// tracksBlob is the JSONB representation of a timeline's track data.
type tracksBlob struct {
	VideoTracks []domain.TrackItem `json:"video_tracks"`
	AudioTracks []domain.TrackItem `json:"audio_tracks"`
	Subtitles   []domain.Subtitle  `json:"subtitles"`
}

func (r *PostgresPostproductionRepository) SaveTimeline(ctx context.Context, t *domain.Timeline) error {
	tracks, _ := json.Marshal(tracksBlob{
		VideoTracks: orEmptyT(t.VideoTracks),
		AudioTracks: orEmptyT(t.AudioTracks),
		Subtitles:   orEmptyS(t.Subtitles),
	})
	status := string(t.Status)
	if status == "" {
		status = string(domain.TimelineDraft)
	}
	_, err := r.q.Exec(ctx, `
		INSERT INTO postproduction.timelines (id, project_id, episode_id, version, tracks, status, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (episode_id, version) DO UPDATE SET
			tracks = EXCLUDED.tracks, status = EXCLUDED.status`,
		t.ID, t.ProjectID, t.EpisodeID, t.Version, tracks, status, t.CreatedAt)
	return err
}

func orEmptyT(v []domain.TrackItem) []domain.TrackItem {
	if v == nil {
		return []domain.TrackItem{}
	}
	return v
}

func orEmptyS(v []domain.Subtitle) []domain.Subtitle {
	if v == nil {
		return []domain.Subtitle{}
	}
	return v
}

func scanTimeline(sc interface {
	Scan(dest ...interface{}) error
}) (*domain.Timeline, error) {
	var t domain.Timeline
	var tracks []byte
	var status string
	err := sc.Scan(&t.ID, &t.ProjectID, &t.EpisodeID, &t.Version, &tracks, &status, &t.CreatedAt)
	t.Status = domain.TimelineStatus(status)
	var blob tracksBlob
	if json.Unmarshal(tracks, &blob) == nil {
		t.VideoTracks = blob.VideoTracks
		t.AudioTracks = blob.AudioTracks
		t.Subtitles = blob.Subtitles
	}
	return &t, err
}

const timelineColumns = `id, project_id, episode_id, version, tracks, status, created_at`

func (r *PostgresPostproductionRepository) FindTimelineByEpisode(ctx context.Context, episodeID string) (*domain.Timeline, error) {
	t, err := scanTimeline(r.q.QueryRow(ctx, `
		SELECT id, project_id, episode_id, version, tracks, status, created_at FROM postproduction.timelines
		WHERE episode_id = $1 ORDER BY version DESC LIMIT 1`, episodeID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrTimelineNotFound
	}
	return t, err
}

func (r *PostgresPostproductionRepository) ListTimelineVersions(ctx context.Context, episodeID string) ([]*domain.Timeline, error) {
	rows, err := r.q.Query(ctx, `
		SELECT id, project_id, episode_id, version, tracks, status, created_at FROM postproduction.timelines
		WHERE episode_id = $1 ORDER BY version`, episodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*domain.Timeline{}
	for rows.Next() {
		t, err := scanTimeline(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

const renderColumns = `id, timeline_id, project_id, episode_id, object_key, url, status, error, created_at, finished_at`

func scanRender(sc interface {
	Scan(dest ...interface{}) error
}) (*domain.RenderTask, error) {
	var rt domain.RenderTask
	var status string
	err := sc.Scan(&rt.ID, &rt.TimelineID, &rt.ProjectID, &rt.EpisodeID, &rt.ObjectKey,
		&rt.OutputURL, &status, &rt.Error, &rt.CreatedAt, &rt.FinishedAt)
	rt.Status = domain.RenderTaskStatus(status)
	return &rt, err
}

func (r *PostgresPostproductionRepository) SaveRender(ctx context.Context, rt *domain.RenderTask) error {
	status := string(rt.Status)
	if status == "" {
		status = string(domain.RenderQueued)
	}
	_, err := r.q.Exec(ctx, `
		INSERT INTO postproduction.renders (id, timeline_id, project_id, episode_id, object_key, url, status, error, created_at, finished_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		ON CONFLICT (id) DO UPDATE SET
			object_key = EXCLUDED.object_key, url = EXCLUDED.url,
			status = EXCLUDED.status, error = EXCLUDED.error,
			finished_at = EXCLUDED.finished_at`,
		rt.ID, rt.TimelineID, rt.ProjectID, rt.EpisodeID, rt.ObjectKey, rt.OutputURL,
		status, rt.Error, rt.CreatedAt, rt.FinishedAt)
	return err
}

func (r *PostgresPostproductionRepository) FindRenderByID(ctx context.Context, id string) (*domain.RenderTask, error) {
	rt, err := scanRender(r.q.QueryRow(ctx, `SELECT id, timeline_id, project_id, episode_id, object_key, url, status, error, created_at, finished_at FROM postproduction.renders WHERE id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrRenderNotFound
	}
	return rt, err
}

func (r *PostgresPostproductionRepository) ListRendersByEpisode(ctx context.Context, episodeID string) ([]*domain.RenderTask, error) {
	rows, err := r.q.Query(ctx, `SELECT id, timeline_id, project_id, episode_id, object_key, url, status, error, created_at, finished_at FROM postproduction.renders WHERE episode_id = $1 ORDER BY created_at DESC`, episodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*domain.RenderTask{}
	for rows.Next() {
		rt, err := scanRender(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, rt)
	}
	return out, rows.Err()
}
