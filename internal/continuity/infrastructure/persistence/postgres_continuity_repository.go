package persistence

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"

	"dramastudio/internal/continuity/domain"
	"dramastudio/internal/platform/database/postgres"
)

type PostgresContinuityRepository struct {
	q postgres.Querier
}

func NewPostgresContinuityRepository(q postgres.Querier) *PostgresContinuityRepository {
	return &PostgresContinuityRepository{q: q}
}

// --- Checks ---

func (r *PostgresContinuityRepository) SaveCheck(ctx context.Context, c *domain.ContinuityCheck) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO continuity.checks (id, project_id, episode_id, check_type, status, created_at, finished_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status, finished_at = EXCLUDED.finished_at`,
		c.ID, c.ProjectID, c.EpisodeID, string(c.CheckType), string(c.Status), c.CreatedAt, c.FinishedAt)
	return err
}

func (r *PostgresContinuityRepository) FindCheckByID(ctx context.Context, id string) (*domain.ContinuityCheck, error) {
	var c domain.ContinuityCheck
	err := r.q.QueryRow(ctx, `
		SELECT id, project_id, episode_id, check_type, status, created_at, finished_at
		FROM continuity.checks WHERE id = $1`, id).
		Scan(&c.ID, &c.ProjectID, &c.EpisodeID, &c.CheckType, &c.Status, &c.CreatedAt, &c.FinishedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrCheckNotFound
	}
	return &c, err
}

func (r *PostgresContinuityRepository) ListChecks(ctx context.Context, projectID, episodeID string) ([]*domain.ContinuityCheck, error) {
	sql := `SELECT id, project_id, episode_id, check_type, status, created_at, finished_at FROM continuity.checks WHERE project_id = $1`
	args := []interface{}{projectID}
	if episodeID != "" {
		sql += ` AND episode_id = $2`
		args = append(args, episodeID)
	}
	sql += ` ORDER BY created_at DESC`
	rows, err := r.q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	checks := []*domain.ContinuityCheck{}
	for rows.Next() {
		var c domain.ContinuityCheck
		if err := rows.Scan(&c.ID, &c.ProjectID, &c.EpisodeID, &c.CheckType, &c.Status, &c.CreatedAt, &c.FinishedAt); err != nil {
			return nil, err
		}
		checks = append(checks, &c)
	}
	return checks, rows.Err()
}

// --- Issues ---

func (r *PostgresContinuityRepository) SaveIssue(ctx context.Context, i *domain.ContinuityIssue) error {
	status := string(i.Status)
	if status == "" {
		status = string(domain.IssueOpen)
	}
	_, err := r.q.Exec(ctx, `
		INSERT INTO continuity.issues (id, check_id, project_id, episode_id, scene_id, category, severity, entity, evidence, expected_state, actual_state, cause, resolution, status, created_at, resolved_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
		ON CONFLICT (id) DO UPDATE SET
			resolution = EXCLUDED.resolution, status = EXCLUDED.status,
			resolved_at = EXCLUDED.resolved_at`,
		i.ID, i.CheckID, i.ProjectID, i.EpisodeID, i.SceneID, i.Category,
		string(i.Severity), i.Entity, i.Evidence, i.ExpectedState, i.ActualState,
		i.Cause, i.Resolution, status, i.CreatedAt, i.ResolvedAt)
	return err
}

const issueColumns = `id, check_id, project_id, episode_id, scene_id, category, severity, entity, evidence, expected_state, actual_state, cause, resolution, status, created_at, resolved_at`

func scanIssue(sc interface {
	Scan(dest ...interface{}) error
}) (*domain.ContinuityIssue, error) {
	var i domain.ContinuityIssue
	var status string
	err := sc.Scan(&i.ID, &i.CheckID, &i.ProjectID, &i.EpisodeID, &i.SceneID,
		&i.Category, &i.Severity, &i.Entity, &i.Evidence, &i.ExpectedState,
		&i.ActualState, &i.Cause, &i.Resolution, &status, &i.CreatedAt, &i.ResolvedAt)
	i.Status = domain.IssueStatus(status)
	return &i, err
}

func (r *PostgresContinuityRepository) FindIssueByID(ctx context.Context, id string) (*domain.ContinuityIssue, error) {
	i, err := scanIssue(r.q.QueryRow(ctx, `SELECT id, check_id, project_id, episode_id, scene_id, category, severity, entity, evidence, expected_state, actual_state, cause, resolution, status, created_at, resolved_at FROM continuity.issues WHERE id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrIssueNotFound
	}
	return i, err
}

func (r *PostgresContinuityRepository) ListIssues(ctx context.Context, projectID string, f domain.IssueFilter) ([]*domain.ContinuityIssue, error) {
	rows, err := r.q.Query(ctx, `SELECT id, check_id, project_id, episode_id, scene_id, category, severity, entity, evidence, expected_state, actual_state, cause, resolution, status, created_at, resolved_at
		FROM continuity.issues
		WHERE project_id = $1
		  AND ($2::text = '' OR episode_id = $2)
		  AND ($3::text = '' OR status = $3)
		  AND ($4::text = '' OR severity = $4)
		  AND ($5::text = '' OR category = $5)
		ORDER BY created_at DESC`,
		projectID, f.EpisodeID, string(f.Status), string(f.Severity), f.Category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	issues := []*domain.ContinuityIssue{}
	for rows.Next() {
		i, err := scanIssue(rows)
		if err != nil {
			return nil, err
		}
		issues = append(issues, i)
	}
	return issues, rows.Err()
}

// --- Timeline events ---

func (r *PostgresContinuityRepository) SaveTimelineEvent(ctx context.Context, e *domain.TimelineEvent) error {
	participants, _ := json.Marshal(orEmpty(e.Participants))
	_, err := r.q.Exec(ctx, `
		INSERT INTO continuity.timeline_events (id, project_id, episode_id, scene_id, world_time, event_order, participants, location_id, duration_sec, description)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		ON CONFLICT (id) DO UPDATE SET
			world_time = EXCLUDED.world_time, event_order = EXCLUDED.event_order,
			participants = EXCLUDED.participants, location_id = EXCLUDED.location_id,
			duration_sec = EXCLUDED.duration_sec, description = EXCLUDED.description`,
		e.ID, e.ProjectID, e.EpisodeID, e.SceneID, e.WorldTime, e.EventOrder,
		participants, e.LocationID, e.DurationSec, e.Description)
	return err
}

func orEmpty(v []string) []string {
	if v == nil {
		return []string{}
	}
	return v
}

func (r *PostgresContinuityRepository) ListTimelineEvents(ctx context.Context, projectID, episodeID string) ([]*domain.TimelineEvent, error) {
	sql := `SELECT id, project_id, episode_id, scene_id, world_time, event_order, participants, location_id, duration_sec, description
		FROM continuity.timeline_events WHERE project_id = $1`
	args := []interface{}{projectID}
	if episodeID != "" {
		sql += ` AND episode_id = $2`
		args = append(args, episodeID)
	}
	sql += ` ORDER BY event_order, world_time`
	rows, err := r.q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events := []*domain.TimelineEvent{}
	for rows.Next() {
		var e domain.TimelineEvent
		var participants []byte
		if err := rows.Scan(&e.ID, &e.ProjectID, &e.EpisodeID, &e.SceneID, &e.WorldTime,
			&e.EventOrder, &participants, &e.LocationID, &e.DurationSec, &e.Description); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(participants, &e.Participants)
		events = append(events, &e)
	}
	return events, rows.Err()
}
