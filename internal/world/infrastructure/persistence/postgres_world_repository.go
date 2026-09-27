package persistence

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"

	"dramastudio/internal/platform/database/postgres"
	"dramastudio/internal/world/domain"
)

type PostgresWorldRepository struct {
	q postgres.Querier
}

func NewPostgresWorldRepository(q postgres.Querier) *PostgresWorldRepository {
	return &PostgresWorldRepository{q: q}
}

func (r *PostgresWorldRepository) SaveLocation(ctx context.Context, loc *domain.Location) error {
	if _, err := r.q.Exec(ctx, `
		INSERT INTO world.locations (id, project_id, name, kind, description, parent_id, visual_ref)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name, kind = EXCLUDED.kind,
			description = EXCLUDED.description, parent_id = EXCLUDED.parent_id,
			visual_ref = EXCLUDED.visual_ref`,
		loc.ID, loc.ProjectID, loc.Name, loc.Kind, loc.Description, loc.ParentID, loc.VisualRef); err != nil {
		return err
	}
	if _, err := r.q.Exec(ctx, `DELETE FROM world.location_variants WHERE location_id = $1`, loc.ID); err != nil {
		return err
	}
	for _, v := range loc.Variants {
		attrs := mergeAssetURL(mustJSON(v.Attributes), v.AssetURL)
		if _, err := r.q.Exec(ctx, `
			INSERT INTO world.location_variants (id, location_id, name, attributes)
			VALUES ($1,$2,$3,$4)`,
			v.ID, loc.ID, v.Name, attrs); err != nil {
			return err
		}
	}
	return nil
}

func mustJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

func mergeAssetURL(attrs []byte, assetURL string) []byte {
	if assetURL == "" {
		return attrs
	}
	var m map[string]string
	if err := json.Unmarshal(attrs, &m); err != nil || m == nil {
		m = map[string]string{}
	}
	m["asset_url"] = assetURL
	b, _ := json.Marshal(m)
	return b
}

func (r *PostgresWorldRepository) FindLocationByID(ctx context.Context, id string) (*domain.Location, error) {
	var loc domain.Location
	err := r.q.QueryRow(ctx, `
		SELECT id, project_id, name, kind, description, parent_id, visual_ref
		FROM world.locations WHERE id = $1`, id).
		Scan(&loc.ID, &loc.ProjectID, &loc.Name, &loc.Kind, &loc.Description, &loc.ParentID, &loc.VisualRef)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrLocationNotFound
	}
	if err != nil {
		return nil, err
	}
	rows, err := r.q.Query(ctx, `
		SELECT id, name, attributes FROM world.location_variants WHERE location_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var v domain.LocationVariant
		var attrs []byte
		if err := rows.Scan(&v.ID, &v.Name, &attrs); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(attrs, &v.Attributes)
		if u, ok := v.Attributes["asset_url"]; ok {
			v.AssetURL = u
			delete(v.Attributes, "asset_url")
		}
		loc.Variants = append(loc.Variants, v)
	}
	return &loc, rows.Err()
}

func (r *PostgresWorldRepository) ListLocations(ctx context.Context, projectID string) ([]*domain.Location, error) {
	rows, err := r.q.Query(ctx, `
		SELECT id, project_id, name, kind, description, parent_id, visual_ref
		FROM world.locations WHERE project_id = $1 ORDER BY name`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	locs := []*domain.Location{}
	for rows.Next() {
		var loc domain.Location
		if err := rows.Scan(&loc.ID, &loc.ProjectID, &loc.Name, &loc.Kind,
			&loc.Description, &loc.ParentID, &loc.VisualRef); err != nil {
			return nil, err
		}
		locs = append(locs, &loc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Hydrate variants in a second pass to keep the query surface simple.
	for _, loc := range locs {
		full, err := r.FindLocationByID(ctx, loc.ID)
		if err != nil {
			return nil, err
		}
		loc.Variants = full.Variants
	}
	return locs, nil
}

func (r *PostgresWorldRepository) SaveProp(ctx context.Context, p *domain.Prop) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO world.props (id, project_id, name, description, location_id, visual_ref)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name, description = EXCLUDED.description,
			location_id = EXCLUDED.location_id, visual_ref = EXCLUDED.visual_ref`,
		p.ID, p.ProjectID, p.Name, p.Description, p.LocationID, p.VisualRef)
	return err
}

func (r *PostgresWorldRepository) FindPropByID(ctx context.Context, id string) (*domain.Prop, error) {
	var p domain.Prop
	err := r.q.QueryRow(ctx, `
		SELECT id, project_id, name, description, location_id, visual_ref
		FROM world.props WHERE id = $1`, id).
		Scan(&p.ID, &p.ProjectID, &p.Name, &p.Description, &p.LocationID, &p.VisualRef)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrLocationNotFound
	}
	return &p, err
}

func (r *PostgresWorldRepository) ListProps(ctx context.Context, projectID string) ([]*domain.Prop, error) {
	rows, err := r.q.Query(ctx, `
		SELECT id, project_id, name, description, location_id, visual_ref
		FROM world.props WHERE project_id = $1 ORDER BY name`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	props := []*domain.Prop{}
	for rows.Next() {
		var p domain.Prop
		if err := rows.Scan(&p.ID, &p.ProjectID, &p.Name, &p.Description, &p.LocationID, &p.VisualRef); err != nil {
			return nil, err
		}
		props = append(props, &p)
	}
	return props, rows.Err()
}

func (r *PostgresWorldRepository) SaveWorldRule(ctx context.Context, rule *domain.WorldRule) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO world.world_rules (id, project_id, rule_text, category)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (id) DO UPDATE SET rule_text = EXCLUDED.rule_text, category = EXCLUDED.category`,
		rule.ID, rule.ProjectID, rule.Text, rule.Category)
	return err
}

func (r *PostgresWorldRepository) ListWorldRules(ctx context.Context, projectID string) ([]*domain.WorldRule, error) {
	rows, err := r.q.Query(ctx, `
		SELECT id, project_id, rule_text, category
		FROM world.world_rules WHERE project_id = $1 ORDER BY id`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rules := []*domain.WorldRule{}
	for rows.Next() {
		var rl domain.WorldRule
		if err := rows.Scan(&rl.ID, &rl.ProjectID, &rl.Text, &rl.Category); err != nil {
			return nil, err
		}
		rules = append(rules, &rl)
	}
	return rules, rows.Err()
}
