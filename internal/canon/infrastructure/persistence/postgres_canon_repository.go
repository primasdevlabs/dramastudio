package persistence

import (
	"context"
	"encoding/json"

	"dramastudio/internal/canon/domain"
	"dramastudio/internal/platform/database/postgres"
)

type PostgresCanonRepository struct {
	q postgres.Querier
}

func NewPostgresCanonRepository(q postgres.Querier) *PostgresCanonRepository {
	return &PostgresCanonRepository{q: q}
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

const factCols = `id, project_id, entity_id, subject, predicate, object, fact_type,
	introduced_episode, effective_from, effective_until, source, confidence, status, version, created_at`

func (r *PostgresCanonRepository) SaveFact(ctx context.Context, f *domain.StoryFact) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO canon.story_facts (id, project_id, entity_id, subject, predicate, object, fact_type,
	introduced_episode, effective_from, effective_until, source, confidence, status, version, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		ON CONFLICT (id) DO UPDATE SET
			entity_id = EXCLUDED.entity_id, subject = EXCLUDED.subject,
			predicate = EXCLUDED.predicate, object = EXCLUDED.object,
			fact_type = EXCLUDED.fact_type, effective_from = EXCLUDED.effective_from,
			effective_until = EXCLUDED.effective_until, source = EXCLUDED.source,
			confidence = EXCLUDED.confidence, status = EXCLUDED.status,
			version = EXCLUDED.version`,
		f.ID, f.ProjectID, f.EntityID, f.Subject, f.Predicate, f.Object, f.Type,
		f.IntroducedEpisode, f.EffectiveFrom, f.EffectiveUntil, f.Source,
		f.Confidence, f.Status, f.Version, f.CreatedAt)
	return err
}

func (r *PostgresCanonRepository) FindFactByID(ctx context.Context, id string) (*domain.StoryFact, error) {
	row := r.q.QueryRow(ctx, `SELECT id, project_id, entity_id, subject, predicate, object, fact_type,
	introduced_episode, effective_from, effective_until, source, confidence, status, version, created_at FROM canon.story_facts WHERE id = $1`, id)
	return scanFact(row)
}

func (r *PostgresCanonRepository) ListFacts(ctx context.Context, projectID string) ([]*domain.StoryFact, error) {
	rows, err := r.q.Query(ctx, `SELECT id, project_id, entity_id, subject, predicate, object, fact_type,
	introduced_episode, effective_from, effective_until, source, confidence, status, version, created_at FROM canon.story_facts WHERE project_id = $1 ORDER BY created_at`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.StoryFact
	for rows.Next() {
		f, err := scanFact(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func (r *PostgresCanonRepository) ListFactsForEntity(ctx context.Context, projectID, entityID string) ([]*domain.StoryFact, error) {
	rows, err := r.q.Query(ctx, `SELECT id, project_id, entity_id, subject, predicate, object, fact_type,
	introduced_episode, effective_from, effective_until, source, confidence, status, version, created_at FROM canon.story_facts WHERE project_id = $1 AND (entity_id = $2 OR subject = $2) ORDER BY created_at`, projectID, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.StoryFact
	for rows.Next() {
		f, err := scanFact(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func scanFact(row rowScanner) (*domain.StoryFact, error) {
	var f domain.StoryFact
	err := row.Scan(&f.ID, &f.ProjectID, &f.EntityID, &f.Subject, &f.Predicate,
		&f.Object, &f.Type, &f.IntroducedEpisode, &f.EffectiveFrom, &f.EffectiveUntil,
		&f.Source, &f.Confidence, &f.Status, &f.Version, &f.CreatedAt)
	if postgres.IsNoRows(err) {
		return nil, domain.ErrFactNotFound
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *PostgresCanonRepository) SaveFactVersion(ctx context.Context, factID string, version int, value []byte) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO canon.fact_versions (fact_id, version, value)
		VALUES ($1,$2,$3) ON CONFLICT (fact_id, version) DO UPDATE SET value = EXCLUDED.value`,
		factID, version, value)
	return err
}

func (r *PostgresCanonRepository) ListFactVersions(ctx context.Context, factID string) ([]*domain.FactVersion, error) {
	rows, err := r.q.Query(ctx, `
		SELECT fact_id, version, value, updated_at
		FROM canon.fact_versions WHERE fact_id = $1 ORDER BY version`, factID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]*domain.FactVersion, 0)
	for rows.Next() {
		v := &domain.FactVersion{}
		if err := rows.Scan(&v.FactID, &v.Version, &v.Value, &v.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (r *PostgresCanonRepository) GetKnowledgeState(ctx context.Context, characterID, episodeID string) (*domain.KnowledgeState, error) {
	var ks domain.KnowledgeState
	var ids []byte
	err := r.q.QueryRow(ctx, `
		SELECT character_id, episode_id, known_fact_ids
		FROM canon.knowledge_states WHERE character_id = $1 AND episode_id = $2`,
		characterID, episodeID).Scan(&ks.CharacterID, &ks.EpisodeID, &ids)
	if postgres.IsNoRows(err) {
		return nil, domain.ErrKnowledgeMissing
	}
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(ids, &ks.KnownFactIDs)
	return &ks, nil
}

func (r *PostgresCanonRepository) SaveKnowledgeState(ctx context.Context, ks *domain.KnowledgeState) error {
	ids, err := json.Marshal(ks.KnownFactIDs)
	if err != nil {
		return err
	}
	_, err = r.q.Exec(ctx, `
		INSERT INTO canon.knowledge_states (character_id, episode_id, known_fact_ids, updated_at)
		VALUES ($1,$2,$3, now())
		ON CONFLICT (character_id, episode_id) DO UPDATE SET
			known_fact_ids = EXCLUDED.known_fact_ids, updated_at = now()`,
		ks.CharacterID, ks.EpisodeID, ids)
	return err
}

func (r *PostgresCanonRepository) SaveRule(ctx context.Context, rule *domain.CanonRule) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO canon.canon_rules (id, project_id, rule_text, enforced)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (id) DO UPDATE SET rule_text = EXCLUDED.rule_text, enforced = EXCLUDED.enforced`,
		rule.ID, rule.ProjectID, rule.RuleText, rule.Enforced)
	return err
}

func (r *PostgresCanonRepository) ListRules(ctx context.Context, projectID string) ([]*domain.CanonRule, error) {
	rows, err := r.q.Query(ctx, `SELECT id, project_id, rule_text, enforced FROM canon.canon_rules WHERE project_id = $1`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.CanonRule
	for rows.Next() {
		var rule domain.CanonRule
		if err := rows.Scan(&rule.ID, &rule.ProjectID, &rule.RuleText, &rule.Enforced); err != nil {
			return nil, err
		}
		out = append(out, &rule)
	}
	return out, rows.Err()
}
