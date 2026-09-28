package persistence

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"

	"dramastudio/internal/media/domain"
	"dramastudio/internal/platform/database/postgres"
)

type PostgresMediaRepository struct {
	q postgres.Querier
}

func NewPostgresMediaRepository(q postgres.Querier) *PostgresMediaRepository {
	return &PostgresMediaRepository{q: q}
}

func mustJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

// --- Assets ---

func (r *PostgresMediaRepository) SaveAsset(ctx context.Context, a *domain.Asset) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO media.assets (id, project_id, type, character_id, location_id, episode_id, scene_id, shot_id,
			provider, model, prompt, spec, reference_assets, parameters, status, cost, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status, cost = EXCLUDED.cost,
			provider = EXCLUDED.provider, model = EXCLUDED.model`,
		a.ID, a.ProjectID, string(a.Type), a.CharacterID, a.LocationID, a.EpisodeID,
		a.SceneID, a.ShotID, a.Provider, a.Model, a.Prompt, mustJSON(a.Spec),
		mustJSON(orEmptyStrs(a.ReferenceAssets)), mustJSON(a.Parameters),
		string(a.Status), a.Cost, a.CreatedAt)
	return err
}

func orEmptyStrs(v []string) []string {
	if v == nil {
		return []string{}
	}
	return v
}

const assetColumns = `id, project_id, type, character_id, location_id, episode_id, scene_id, shot_id,
	provider, model, prompt, spec, reference_assets, parameters, status, cost, created_at`

func scanAsset(sc interface {
	Scan(dest ...interface{}) error
}) (*domain.Asset, error) {
	var a domain.Asset
	var spec, refs, params []byte
	err := sc.Scan(&a.ID, &a.ProjectID, &a.Type, &a.CharacterID, &a.LocationID,
		&a.EpisodeID, &a.SceneID, &a.ShotID, &a.Provider, &a.Model, &a.Prompt,
		&spec, &refs, &params, &a.Status, &a.Cost, &a.CreatedAt)
	_ = json.Unmarshal(spec, &a.Spec)
	_ = json.Unmarshal(refs, &a.ReferenceAssets)
	_ = json.Unmarshal(params, &a.Parameters)
	return &a, err
}

func (r *PostgresMediaRepository) FindAssetByID(ctx context.Context, id string) (*domain.Asset, error) {
	a, err := scanAsset(r.q.QueryRow(ctx, `SELECT id, project_id, type, character_id, location_id, episode_id, scene_id, shot_id,
	provider, model, prompt, spec, reference_assets, parameters, status, cost, created_at FROM media.assets WHERE id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrAssetNotFound
	}
	return a, err
}

func (r *PostgresMediaRepository) ListAssetsByProject(ctx context.Context, projectID string) ([]*domain.Asset, error) {
	rows, err := r.q.Query(ctx, `SELECT id, project_id, type, character_id, location_id, episode_id, scene_id, shot_id,
	provider, model, prompt, spec, reference_assets, parameters, status, cost, created_at FROM media.assets WHERE project_id = $1 ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	assets := []*domain.Asset{}
	for rows.Next() {
		a, err := scanAsset(rows)
		if err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, rows.Err()
}

// --- Versions ---

func (r *PostgresMediaRepository) SaveAssetVersion(ctx context.Context, v *domain.AssetVersion) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO media.asset_versions (id, asset_id, version, object_key, url, status, cost, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (asset_id, version) DO UPDATE SET
			object_key = EXCLUDED.object_key, url = EXCLUDED.url,
			status = EXCLUDED.status, cost = EXCLUDED.cost`,
		v.ID, v.AssetID, v.Version, v.ObjectKey, v.URL, string(v.Status), v.Cost, v.CreatedAt)
	return err
}

func (r *PostgresMediaRepository) FindAssetVersion(ctx context.Context, assetID string, version int) (*domain.AssetVersion, error) {
	var v domain.AssetVersion
	err := r.q.QueryRow(ctx, `
		SELECT id, asset_id, version, object_key, url, status, cost, created_at
		FROM media.asset_versions WHERE asset_id = $1 AND version = $2`, assetID, version).
		Scan(&v.ID, &v.AssetID, &v.Version, &v.ObjectKey, &v.URL, &v.Status, &v.Cost, &v.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrVersionNotFound
	}
	return &v, err
}

func (r *PostgresMediaRepository) ListAssetVersions(ctx context.Context, assetID string) ([]*domain.AssetVersion, error) {
	rows, err := r.q.Query(ctx, `
		SELECT id, asset_id, version, object_key, url, status, cost, created_at
		FROM media.asset_versions WHERE asset_id = $1 ORDER BY version`, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	versions := []*domain.AssetVersion{}
	for rows.Next() {
		var v domain.AssetVersion
		if err := rows.Scan(&v.ID, &v.AssetID, &v.Version, &v.ObjectKey, &v.URL, &v.Status, &v.Cost, &v.CreatedAt); err != nil {
			return nil, err
		}
		versions = append(versions, &v)
	}
	return versions, rows.Err()
}

// --- Generation jobs ---

const jobColumns = `id, project_id, asset_id, capability, provider, model, provider_job_id,
	input, output_url, status, attempt, cost, error, started_at, completed_at`

func scanJob(sc interface {
	Scan(dest ...interface{}) error
}) (*domain.GenerationJob, error) {
	var j domain.GenerationJob
	err := sc.Scan(&j.ID, &j.ProjectID, &j.AssetID, &j.Capability, &j.Provider,
		&j.Model, &j.ProviderJobID, &j.Input, &j.OutputURL, &j.Status,
		&j.Attempt, &j.Cost, &j.Error, &j.StartedAt, &j.CompletedAt)
	return &j, err
}

func (r *PostgresMediaRepository) SaveJob(ctx context.Context, j *domain.GenerationJob) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO media.generation_jobs (id, project_id, asset_id, capability, provider, model, provider_job_id, input, output_url, status, attempt, cost, error, started_at, completed_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		ON CONFLICT (id) DO UPDATE SET
			provider_job_id = EXCLUDED.provider_job_id, output_url = EXCLUDED.output_url,
			status = EXCLUDED.status, attempt = EXCLUDED.attempt, cost = EXCLUDED.cost,
			error = EXCLUDED.error, completed_at = EXCLUDED.completed_at`,
		j.ID, j.ProjectID, j.AssetID, j.Capability, j.Provider, j.Model,
		j.ProviderJobID, j.Input, j.OutputURL, string(j.Status), j.Attempt,
		j.Cost, j.Error, j.StartedAt, j.CompletedAt)
	return err
}

func (r *PostgresMediaRepository) FindJobByID(ctx context.Context, id string) (*domain.GenerationJob, error) {
	j, err := scanJob(r.q.QueryRow(ctx, `SELECT id, project_id, asset_id, capability, provider, model, provider_job_id,
	input, output_url, status, attempt, cost, error, started_at, completed_at FROM media.generation_jobs WHERE id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrJobNotFound
	}
	return j, err
}

func (r *PostgresMediaRepository) FindJobByProviderJobID(ctx context.Context, providerJobID string) (*domain.GenerationJob, error) {
	j, err := scanJob(r.q.QueryRow(ctx, `SELECT id, project_id, asset_id, capability, provider, model, provider_job_id,
	input, output_url, status, attempt, cost, error, started_at, completed_at FROM media.generation_jobs WHERE provider_job_id = $1`, providerJobID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrJobNotFound
	}
	return j, err
}

func (r *PostgresMediaRepository) ListJobsByProject(ctx context.Context, projectID string) ([]*domain.GenerationJob, error) {
	rows, err := r.q.Query(ctx, `SELECT id, project_id, asset_id, capability, provider, model, provider_job_id,
	input, output_url, status, attempt, cost, error, started_at, completed_at FROM media.generation_jobs WHERE project_id = $1 ORDER BY started_at DESC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	jobs := []*domain.GenerationJob{}
	for rows.Next() {
		j, err := scanJob(rows)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, rows.Err()
}
