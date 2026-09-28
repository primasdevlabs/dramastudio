package persistence

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"

	"dramastudio/internal/platform/database/postgres"
	"dramastudio/internal/production/domain"
)

type PostgresProductionRepository struct {
	q postgres.Querier
}

func NewPostgresProductionRepository(q postgres.Querier) *PostgresProductionRepository {
	return &PostgresProductionRepository{q: q}
}

func mustJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

// --- Productions ---

func (r *PostgresProductionRepository) SaveProduction(ctx context.Context, p *domain.Production) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO production.productions (id, project_id)
		VALUES ($1,$2)
		ON CONFLICT (project_id) DO NOTHING`, p.ID, p.ProjectID)
	return err
}

func (r *PostgresProductionRepository) FindProductionByProject(ctx context.Context, projectID string) (*domain.Production, error) {
	var p domain.Production
	err := r.q.QueryRow(ctx, `
		SELECT id, project_id FROM production.productions WHERE project_id = $1`, projectID).
		Scan(&p.ID, &p.ProjectID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrRunNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// --- Runs ---

func (r *PostgresProductionRepository) SaveRun(ctx context.Context, run *domain.ProductionRun) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO production.runs (id, production_id, project_id, episode_id, stage, status, bible_version, workflow_id, snapshot, created_at, completed_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (id) DO UPDATE SET
			stage = EXCLUDED.stage, status = EXCLUDED.status,
			bible_version = EXCLUDED.bible_version, workflow_id = EXCLUDED.workflow_id,
			snapshot = EXCLUDED.snapshot, completed_at = EXCLUDED.completed_at`,
		run.ID, run.ProductionID, run.ProjectID, run.EpisodeID, string(run.Stage),
		string(run.Status), run.BibleVersion, run.WorkflowID, mustJSON(run.Snapshot),
		run.CreatedAt, run.CompletedAt)
	return err
}

func (r *PostgresProductionRepository) FindRunByID(ctx context.Context, id string) (*domain.ProductionRun, error) {
	var run domain.ProductionRun
	var snapshot []byte
	err := r.q.QueryRow(ctx, `
		SELECT id, production_id, project_id, episode_id, stage, status, bible_version, workflow_id, snapshot, created_at, completed_at
		FROM production.runs WHERE id = $1`, id).
		Scan(&run.ID, &run.ProductionID, &run.ProjectID, &run.EpisodeID, &run.Stage,
			&run.Status, &run.BibleVersion, &run.WorkflowID, &snapshot, &run.CreatedAt, &run.CompletedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrRunNotFound
	}
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(snapshot, &run.Snapshot)
	return &run, nil
}

func (r *PostgresProductionRepository) ListRunsByProject(ctx context.Context, projectID string) ([]*domain.ProductionRun, error) {
	rows, err := r.q.Query(ctx, `
		SELECT id, production_id, project_id, episode_id, stage, status, bible_version, workflow_id, snapshot, created_at, completed_at
		FROM production.runs WHERE project_id = $1 ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	runs := []*domain.ProductionRun{}
	for rows.Next() {
		var run domain.ProductionRun
		var snapshot []byte
		if err := rows.Scan(&run.ID, &run.ProductionID, &run.ProjectID, &run.EpisodeID,
			&run.Stage, &run.Status, &run.BibleVersion, &run.WorkflowID, &snapshot,
			&run.CreatedAt, &run.CompletedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(snapshot, &run.Snapshot)
		runs = append(runs, &run)
	}
	return runs, rows.Err()
}

// --- Jobs ---

func (r *PostgresProductionRepository) SaveJob(ctx context.Context, job *domain.ProductionJob) error {
	var key *string
	if job.IdempotencyKey != "" {
		key = &job.IdempotencyKey
	}
	_, err := r.q.Exec(ctx, `
		INSERT INTO production.jobs (id, production_id, project_id, run_id, episode_id, scene_id, shot_id, kind, status, attempt, result_url, error, idempotency_key, created_at, completed_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status, attempt = EXCLUDED.attempt,
			result_url = EXCLUDED.result_url, error = EXCLUDED.error,
			completed_at = EXCLUDED.completed_at`,
		job.ID, job.ProductionID, job.ProjectID, job.RunID, job.EpisodeID, job.SceneID,
		job.ShotID, job.Kind, string(job.Status), job.Attempt, job.ResultURL, job.Error,
		key, job.CreatedAt, job.CompletedAt)
	return err
}

func scanJob(sc interface {
	Scan(dest ...interface{}) error
}) (*domain.ProductionJob, error) {
	var j domain.ProductionJob
	var key *string
	err := sc.Scan(&j.ID, &j.ProductionID, &j.ProjectID, &j.RunID, &j.EpisodeID,
		&j.SceneID, &j.ShotID, &j.Kind, &j.Status, &j.Attempt, &j.ResultURL,
		&j.Error, &key, &j.CreatedAt, &j.CompletedAt)
	if key != nil {
		j.IdempotencyKey = *key
	}
	return &j, err
}

const jobColumns = `id, production_id, project_id, run_id, episode_id, scene_id, shot_id, kind, status, attempt, result_url, error, idempotency_key, created_at, completed_at`

func (r *PostgresProductionRepository) FindJobByID(ctx context.Context, id string) (*domain.ProductionJob, error) {
	j, err := scanJob(r.q.QueryRow(ctx, `SELECT id, production_id, project_id, run_id, episode_id, scene_id, shot_id, kind, status, attempt, result_url, error, idempotency_key, created_at, completed_at FROM production.jobs WHERE id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrJobNotFound
	}
	return j, err
}

func (r *PostgresProductionRepository) FindJobByIdempotencyKey(ctx context.Context, key string) (*domain.ProductionJob, error) {
	j, err := scanJob(r.q.QueryRow(ctx, `SELECT id, production_id, project_id, run_id, episode_id, scene_id, shot_id, kind, status, attempt, result_url, error, idempotency_key, created_at, completed_at FROM production.jobs WHERE idempotency_key = $1`, key))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrJobNotFound
	}
	return j, err
}

func (r *PostgresProductionRepository) ListJobsByProject(ctx context.Context, projectID string) ([]*domain.ProductionJob, error) {
	return r.listJobs(ctx, `SELECT id, production_id, project_id, run_id, episode_id, scene_id, shot_id, kind, status, attempt, result_url, error, idempotency_key, created_at, completed_at FROM production.jobs WHERE project_id = $1 ORDER BY created_at`, projectID)
}

func (r *PostgresProductionRepository) ListJobsByRun(ctx context.Context, runID string) ([]*domain.ProductionJob, error) {
	return r.listJobs(ctx, `SELECT id, production_id, project_id, run_id, episode_id, scene_id, shot_id, kind, status, attempt, result_url, error, idempotency_key, created_at, completed_at FROM production.jobs WHERE run_id = $1 ORDER BY created_at`, runID)
}

func (r *PostgresProductionRepository) listJobs(ctx context.Context, sql, arg string) ([]*domain.ProductionJob, error) {
	rows, err := r.q.Query(ctx, sql, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	jobs := []*domain.ProductionJob{}
	for rows.Next() {
		j, err := scanJob(rows)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, rows.Err()
}

// --- Shots ---

func (r *PostgresProductionRepository) SaveShot(ctx context.Context, s *domain.Shot) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO production.shots (id, project_id, episode_id, scene_id, seq, description, camera, characters, location_id, duration_sec, status, approved_asset)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		ON CONFLICT (id) DO UPDATE SET
			description = EXCLUDED.description, camera = EXCLUDED.camera,
			characters = EXCLUDED.characters, location_id = EXCLUDED.location_id,
			duration_sec = EXCLUDED.duration_sec, status = EXCLUDED.status,
			approved_asset = EXCLUDED.approved_asset`,
		s.ID, s.ProjectID, s.EpisodeID, s.SceneID, s.Seq, s.Description,
		mustJSON(s.Camera), mustJSON(orEmpty(s.Characters)), s.LocationID,
		s.DurationSec, string(s.Status), s.ApprovedAsset)
	return err
}

func orEmpty(v []string) []string {
	if v == nil {
		return []string{}
	}
	return v
}

func scanShot(sc interface {
	Scan(dest ...interface{}) error
}) (*domain.Shot, error) {
	var s domain.Shot
	var camera, chars []byte
	err := sc.Scan(&s.ID, &s.ProjectID, &s.EpisodeID, &s.SceneID, &s.Seq,
		&s.Description, &camera, &chars, &s.LocationID, &s.DurationSec,
		&s.Status, &s.ApprovedAsset)
	_ = json.Unmarshal(camera, &s.Camera)
	_ = json.Unmarshal(chars, &s.Characters)
	return &s, err
}

const shotColumns = `id, project_id, episode_id, scene_id, seq, description, camera, characters, location_id, duration_sec, status, approved_asset`

func (r *PostgresProductionRepository) FindShotByID(ctx context.Context, id string) (*domain.Shot, error) {
	s, err := scanShot(r.q.QueryRow(ctx, `SELECT id, project_id, episode_id, scene_id, seq, description, camera, characters, location_id, duration_sec, status, approved_asset FROM production.shots WHERE id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrShotNotFound
	}
	return s, err
}

func (r *PostgresProductionRepository) ListShots(ctx context.Context, episodeID, sceneID string) ([]*domain.Shot, error) {
	var rows pgx.Rows
	var err error
	if sceneID != "" {
		rows, err = r.q.Query(ctx, `SELECT id, project_id, episode_id, scene_id, seq, description, camera, characters, location_id, duration_sec, status, approved_asset FROM production.shots WHERE scene_id = $1 ORDER BY seq`, sceneID)
	} else {
		rows, err = r.q.Query(ctx, `SELECT id, project_id, episode_id, scene_id, seq, description, camera, characters, location_id, duration_sec, status, approved_asset FROM production.shots WHERE episode_id = $1 ORDER BY scene_id, seq`, episodeID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	shots := []*domain.Shot{}
	for rows.Next() {
		s, err := scanShot(rows)
		if err != nil {
			return nil, err
		}
		shots = append(shots, s)
	}
	return shots, rows.Err()
}

// --- Approvals ---

func (r *PostgresProductionRepository) SaveApproval(ctx context.Context, a *domain.ApprovalRequest) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO production.approvals (id, project_id, episode_id, stage, target_id, decision, notes, decided_by, submitted_at, decided_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		ON CONFLICT (id) DO UPDATE SET
			decision = EXCLUDED.decision, notes = EXCLUDED.notes,
			decided_by = EXCLUDED.decided_by, decided_at = EXCLUDED.decided_at`,
		a.ID, a.ProjectID, a.EpisodeID, a.Stage, a.TargetID, string(a.Decision),
		a.Notes, a.DecidedBy, a.SubmittedAt, a.DecidedAt)
	return err
}

func scanApproval(sc interface {
	Scan(dest ...interface{}) error
}) (*domain.ApprovalRequest, error) {
	var a domain.ApprovalRequest
	err := sc.Scan(&a.ID, &a.ProjectID, &a.EpisodeID, &a.Stage, &a.TargetID,
		&a.Decision, &a.Notes, &a.DecidedBy, &a.SubmittedAt, &a.DecidedAt)
	return &a, err
}

const approvalColumns = `id, project_id, episode_id, stage, target_id, decision, notes, decided_by, submitted_at, decided_at`

func (r *PostgresProductionRepository) FindApprovalByID(ctx context.Context, id string) (*domain.ApprovalRequest, error) {
	a, err := scanApproval(r.q.QueryRow(ctx, `SELECT id, project_id, episode_id, stage, target_id, decision, notes, decided_by, submitted_at, decided_at FROM production.approvals WHERE id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrApprovalNotFound
	}
	return a, err
}

func (r *PostgresProductionRepository) ListApprovalsByProject(ctx context.Context, projectID string, pendingOnly bool) ([]*domain.ApprovalRequest, error) {
	sql := `SELECT id, project_id, episode_id, stage, target_id, decision, notes, decided_by, submitted_at, decided_at FROM production.approvals WHERE project_id = $1`
	if pendingOnly {
		sql += ` AND decided_at IS NULL`
	}
	sql += ` ORDER BY submitted_at DESC`
	rows, err := r.q.Query(ctx, sql, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	apps := []*domain.ApprovalRequest{}
	for rows.Next() {
		a, err := scanApproval(rows)
		if err != nil {
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, rows.Err()
}
