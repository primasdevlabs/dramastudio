package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"dramastudio/internal/platform/database/postgres"
	"dramastudio/internal/publishing/domain"
)

type PostgresPublishingRepository struct {
	q postgres.Querier
}

func NewPostgresPublishingRepository(q postgres.Querier) *PostgresPublishingRepository {
	return &PostgresPublishingRepository{q: q}
}

func mustJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

// --- Channels ---

func (r *PostgresPublishingRepository) SaveChannel(ctx context.Context, c *domain.Channel) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO publishing.channels (id, project_id, platform, account_ref, config, enabled, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (id) DO UPDATE SET
			account_ref = EXCLUDED.account_ref, config = EXCLUDED.config,
			enabled = EXCLUDED.enabled`,
		c.ID, c.ProjectID, string(c.Platform), c.AccountRef, mustJSON(c.Config), c.Enabled, c.CreatedAt)
	return err
}

func scanChannel(sc interface {
	Scan(dest ...interface{}) error
}) (*domain.Channel, error) {
	var c domain.Channel
	var cfg []byte
	err := sc.Scan(&c.ID, &c.ProjectID, &c.Platform, &c.AccountRef, &cfg, &c.Enabled, &c.CreatedAt)
	_ = json.Unmarshal(cfg, &c.Config)
	return &c, err
}

func (r *PostgresPublishingRepository) FindChannelByID(ctx context.Context, id string) (*domain.Channel, error) {
	c, err := scanChannel(r.q.QueryRow(ctx, `
		SELECT id, project_id, platform, account_ref, config, enabled, created_at
		FROM publishing.channels WHERE id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrChannelNotFound
	}
	return c, err
}

func (r *PostgresPublishingRepository) ListChannels(ctx context.Context, projectID string) ([]*domain.Channel, error) {
	rows, err := r.q.Query(ctx, `
		SELECT id, project_id, platform, account_ref, config, enabled, created_at
		FROM publishing.channels WHERE project_id = $1 ORDER BY created_at`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*domain.Channel{}
	for rows.Next() {
		c, err := scanChannel(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// --- Publications ---

const (
	pubColumns = `id, project_id, episode_id, channel_id, metadata, video_url, scheduled_at, status, platform_response, external_id, published_at, idempotency_key, created_at`

	queryFindPublicationByID = `SELECT id, project_id, episode_id, channel_id, metadata, video_url, scheduled_at, status, platform_response, external_id, published_at, idempotency_key, created_at FROM publishing.publications WHERE id = $1`
)

func scanPublication(sc interface {
	Scan(dest ...interface{}) error
}) (*domain.Publication, error) {
	var p domain.Publication
	var meta, resp []byte
	var status string
	var key *string
	err := sc.Scan(&p.ID, &p.ProjectID, &p.EpisodeID, &p.ChannelID, &meta,
		&p.VideoURL, &p.ScheduledAt, &status, &resp, &p.ExternalID, &p.PublishedAt, &key, &p.CreatedAt)
	p.Status = domain.PublicationStatus(status)
	if key != nil {
		p.IdempotencyKey = *key
	}
	_ = json.Unmarshal(meta, &p.Metadata)
	_ = json.Unmarshal(resp, &p.PlatformResponse)
	return &p, err
}

func (r *PostgresPublishingRepository) SavePublication(ctx context.Context, p *domain.Publication) error {
	var key *string
	if p.IdempotencyKey != "" {
		key = &p.IdempotencyKey
	}
	_, err := r.q.Exec(ctx, `
		INSERT INTO publishing.publications (id, project_id, episode_id, channel_id, metadata, video_url, scheduled_at, status, platform_response, external_id, published_at, idempotency_key, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		ON CONFLICT (id) DO UPDATE SET
			metadata = EXCLUDED.metadata, video_url = EXCLUDED.video_url,
			scheduled_at = EXCLUDED.scheduled_at,
			status = EXCLUDED.status, platform_response = EXCLUDED.platform_response,
			external_id = EXCLUDED.external_id, published_at = EXCLUDED.published_at`,
		p.ID, p.ProjectID, p.EpisodeID, p.ChannelID, mustJSON(p.Metadata),
		p.VideoURL, p.ScheduledAt, string(p.Status), mustJSON(p.PlatformResponse),
		p.ExternalID, p.PublishedAt, key, p.CreatedAt)
	return err
}

func (r *PostgresPublishingRepository) FindPublicationByID(ctx context.Context, id string) (*domain.Publication, error) {
	p, err := scanPublication(r.q.QueryRow(ctx, queryFindPublicationByID, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrPublicationNotFound
	}
	return p, err
}

func (r *PostgresPublishingRepository) FindPublicationByIdempotencyKey(ctx context.Context, key string) (*domain.Publication, error) {
	p, err := scanPublication(r.q.QueryRow(ctx, `SELECT id, project_id, episode_id, channel_id, metadata, video_url, scheduled_at, status, platform_response, external_id, published_at, idempotency_key, created_at FROM publishing.publications WHERE idempotency_key = $1`, key))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrPublicationNotFound
	}
	return p, err
}

func (r *PostgresPublishingRepository) ListPublicationsByProject(ctx context.Context, projectID string) ([]*domain.Publication, error) {
	rows, err := r.q.Query(ctx, `SELECT id, project_id, episode_id, channel_id, metadata, video_url, scheduled_at, status, platform_response, external_id, published_at, idempotency_key, created_at FROM publishing.publications WHERE project_id = $1 ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*domain.Publication{}
	for rows.Next() {
		p, err := scanPublication(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *PostgresPublishingRepository) ListDueScheduled(ctx context.Context, now time.Time) ([]*domain.Publication, error) {
	rows, err := r.q.Query(ctx, `SELECT id, project_id, episode_id, channel_id, metadata, video_url, scheduled_at, status, platform_response, external_id, published_at, idempotency_key, created_at FROM publishing.publications
		WHERE status = 'scheduled' AND scheduled_at IS NOT NULL AND scheduled_at <= $1
		ORDER BY scheduled_at`, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*domain.Publication{}
	for rows.Next() {
		p, err := scanPublication(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
