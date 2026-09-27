package persistence

import (
	"context"
	"encoding/json"
	"time"

	"dramastudio/internal/platform/database/postgres"
	"dramastudio/internal/projects/domain"
)

type PostgresProjectRepository struct {
	q postgres.Querier
}

func NewPostgresProjectRepository(q postgres.Querier) *PostgresProjectRepository {
	return &PostgresProjectRepository{q: q}
}

func (r *PostgresProjectRepository) Save(ctx context.Context, p *domain.Project) error {
	settings, err := json.Marshal(p.Settings)
	if err != nil {
		return err
	}
	policy, err := json.Marshal(p.Policy)
	if err != nil {
		return err
	}
	p.UpdatedAt = time.Now().UTC()
	_, err = r.q.Exec(ctx, `
		INSERT INTO projects.projects (id, org_id, name, description, genre, language, mode, status, settings, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9::jsonb || jsonb_build_object('policy',$10::jsonb,'budget',$11::jsonb),$12,$13)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name, description = EXCLUDED.description,
			genre = EXCLUDED.genre, language = EXCLUDED.language,
			mode = EXCLUDED.mode, status = EXCLUDED.status,
			settings = EXCLUDED.settings, updated_at = EXCLUDED.updated_at`,
		p.ID, p.OrgID, p.Name, p.Description, p.Genre, p.Language, p.Mode, p.Status,
		settings, policy, mustJSON(p.Budget), p.CreatedAt, p.UpdatedAt)
	return err
}

func mustJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

func (r *PostgresProjectRepository) FindByID(ctx context.Context, id domain.ProjectID) (*domain.Project, error) {
	row := r.q.QueryRow(ctx, `
		SELECT id, org_id, name, description, genre, language, mode, status,
		       settings, created_at, updated_at
		FROM projects.projects WHERE id = $1`, id)
	return scanProject(row)
}

func (r *PostgresProjectRepository) ListAll(ctx context.Context, orgID string) ([]*domain.Project, error) {
	rows, err := r.q.Query(ctx, `
		SELECT id, org_id, name, description, genre, language, mode, status,
		       settings, created_at, updated_at
		FROM projects.projects WHERE $1 = '' OR org_id = $1
		ORDER BY created_at`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Project
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanProject(row rowScanner) (*domain.Project, error) {
	var p domain.Project
	var settings []byte
	err := row.Scan(&p.ID, &p.OrgID, &p.Name, &p.Description, &p.Genre, &p.Language,
		&p.Mode, &p.Status, &settings, &p.CreatedAt, &p.UpdatedAt)
	if postgres.IsNoRows(err) {
		return nil, domain.ErrProjectNotFound
	}
	if err != nil {
		return nil, err
	}
	var blob struct {
		domain.Settings
		Policy domain.ProductionPolicy `json:"policy"`
		Budget domain.Budget           `json:"budget"`
	}
	if err := json.Unmarshal(settings, &blob); err == nil {
		p.Settings = blob.Settings
		p.Policy = blob.Policy
		p.Budget = blob.Budget
	}
	return &p, nil
}

func (r *PostgresProjectRepository) SaveBible(ctx context.Context, b *domain.SeriesBible) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO projects.series_bibles
			(id, project_id, version, premise, genre, themes, tone, world_rules,
			 narrative_rules, visual_style, dialogue_style, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		ON CONFLICT (project_id, version) DO UPDATE SET
			premise = EXCLUDED.premise, genre = EXCLUDED.genre,
			themes = EXCLUDED.themes, tone = EXCLUDED.tone,
			world_rules = EXCLUDED.world_rules, narrative_rules = EXCLUDED.narrative_rules,
			visual_style = EXCLUDED.visual_style, dialogue_style = EXCLUDED.dialogue_style`,
		b.ID, b.ProjectID, b.Version, b.Premise, b.Genre, mustJSON(b.Themes), b.Tone,
		mustJSON(b.WorldRules), mustJSON(b.NarrativeRules), b.VisualDirection,
		b.DialogueStyle, b.CreatedAt)
	return err
}

func (r *PostgresProjectRepository) GetLatestBible(ctx context.Context, projectID domain.ProjectID) (*domain.SeriesBible, error) {
	row := r.q.QueryRow(ctx, `
		SELECT id, project_id, version, premise, genre, themes, tone, world_rules,
		       narrative_rules, visual_style, dialogue_style, created_at
		FROM projects.series_bibles
		WHERE project_id = $1 ORDER BY version DESC LIMIT 1`, projectID)
	return scanBible(row)
}

func (r *PostgresProjectRepository) GetBibleVersion(ctx context.Context, projectID domain.ProjectID, version int) (*domain.SeriesBible, error) {
	row := r.q.QueryRow(ctx, `
		SELECT id, project_id, version, premise, genre, themes, tone, world_rules,
		       narrative_rules, visual_style, dialogue_style, created_at
		FROM projects.series_bibles
		WHERE project_id = $1 AND version = $2`, projectID, version)
	return scanBible(row)
}

func (r *PostgresProjectRepository) ListBibleVersions(ctx context.Context, projectID domain.ProjectID) ([]*domain.SeriesBible, error) {
	rows, err := r.q.Query(ctx, `
		SELECT id, project_id, version, premise, genre, themes, tone, world_rules,
		       narrative_rules, visual_style, dialogue_style, created_at
		FROM projects.series_bibles WHERE project_id = $1 ORDER BY version`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.SeriesBible
	for rows.Next() {
		b, err := scanBible(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func scanBible(row rowScanner) (*domain.SeriesBible, error) {
	var b domain.SeriesBible
	var themes, worldRules, narrativeRules []byte
	err := row.Scan(&b.ID, &b.ProjectID, &b.Version, &b.Premise, &b.Genre,
		&themes, &b.Tone, &worldRules, &narrativeRules,
		&b.VisualDirection, &b.DialogueStyle, &b.CreatedAt)
	if postgres.IsNoRows(err) {
		return nil, domain.ErrBibleNotFound
	}
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(themes, &b.Themes)
	_ = json.Unmarshal(worldRules, &b.WorldRules)
	_ = json.Unmarshal(narrativeRules, &b.NarrativeRules)
	return &b, nil
}
