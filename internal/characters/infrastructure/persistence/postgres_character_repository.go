package persistence

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"dramastudio/internal/characters/domain"
	"dramastudio/internal/platform/database/postgres"
)

type PostgresCharacterRepository struct {
	q postgres.Querier
}

func NewPostgresCharacterRepository(q postgres.Querier) *PostgresCharacterRepository {
	return &PostgresCharacterRepository{q: q}
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func mustJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

func (r *PostgresCharacterRepository) Save(ctx context.Context, c *domain.Character) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO characters.characters (id, project_id, name, role, summary, locked)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name, role = EXCLUDED.role,
			summary = EXCLUDED.summary, locked = EXCLUDED.locked`,
		c.ID, c.ProjectID, c.Name, c.Role, c.Bio, c.IsLocked)
	if err != nil {
		return err
	}
	// Snapshot the mutable profile as an immutable version.
	ver := &domain.CharacterVersion{
		ID:           "cver_" + uuid.NewString(),
		CharacterID:  c.ID,
		Version:      c.Version,
		Appearance:   c.Appearance,
		Personality:  c.Personality,
		VoiceProfile: c.VoiceProfile,
		Wardrobe:     c.Wardrobe,
	}
	return r.SaveVersion(ctx, ver)
}

func (r *PostgresCharacterRepository) FindByID(ctx context.Context, id domain.CharacterID) (*domain.Character, error) {
	row := r.q.QueryRow(ctx, `
		SELECT c.id, c.project_id, c.name, c.role, c.summary, c.locked,
		       v.version, v.appearance, v.personality, v.voice_profile, v.wardrobe
		FROM characters.characters c
		LEFT JOIN characters.character_versions v
			ON v.character_id = c.id AND v.version = (
				SELECT MAX(version) FROM characters.character_versions WHERE character_id = c.id)
		WHERE c.id = $1`, id)
	c, err := scanCharacter(row)
	if err != nil {
		return nil, err
	}
	// Relationships live in their own table — load them so reads see the
	// same aggregate the service assembled.
	rels, err := r.ListRelationships(ctx, id)
	if err != nil {
		return nil, err
	}
	c.Relationships = rels
	return c, nil
}

func (r *PostgresCharacterRepository) List(ctx context.Context, projectID string) ([]*domain.Character, error) {
	rows, err := r.q.Query(ctx, `
		SELECT c.id, c.project_id, c.name, c.role, c.summary, c.locked,
		       v.version, v.appearance, v.personality, v.voice_profile, v.wardrobe
		FROM characters.characters c
		LEFT JOIN characters.character_versions v
			ON v.character_id = c.id AND v.version = (
				SELECT MAX(version) FROM characters.character_versions WHERE character_id = c.id)
		WHERE $1 = '' OR c.project_id = $1
		ORDER BY c.created_at`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Character
	for rows.Next() {
		c, err := scanCharacter(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func scanCharacter(row rowScanner) (*domain.Character, error) {
	var c domain.Character
	var version *int
	var appearance, personality, voice, wardrobe []byte
	err := row.Scan(&c.ID, &c.ProjectID, &c.Name, &c.Role, &c.Bio, &c.IsLocked,
		&version, &appearance, &personality, &voice, &wardrobe)
	if postgres.IsNoRows(err) {
		return nil, domain.ErrCharacterNotFound
	}
	if err != nil {
		return nil, err
	}
	if version != nil {
		c.Version = *version
	}
	_ = json.Unmarshal(appearance, &c.Appearance)
	_ = json.Unmarshal(personality, &c.Personality)
	_ = json.Unmarshal(voice, &c.VoiceProfile)
	_ = json.Unmarshal(wardrobe, &c.Wardrobe)
	return &c, nil
}

func (r *PostgresCharacterRepository) SetLocked(ctx context.Context, id domain.CharacterID, locked bool) error {
	tag, err := r.q.Exec(ctx, `UPDATE characters.characters SET locked = $2 WHERE id = $1`, id, locked)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrCharacterNotFound
	}
	return nil
}

func (r *PostgresCharacterRepository) SaveVersion(ctx context.Context, v *domain.CharacterVersion) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO characters.character_versions
			(id, character_id, version, appearance, personality, voice_profile, wardrobe)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (character_id, version) DO NOTHING`,
		v.ID, v.CharacterID, v.Version, mustJSON(v.Appearance), mustJSON(v.Personality),
		mustJSON(v.VoiceProfile), mustJSON(v.Wardrobe))
	return err
}

func (r *PostgresCharacterRepository) FindVersion(ctx context.Context, characterID domain.CharacterID, version int) (*domain.CharacterVersion, error) {
	row := r.q.QueryRow(ctx, `
		SELECT id, character_id, version, appearance, personality, voice_profile, wardrobe
		FROM characters.character_versions WHERE character_id = $1 AND version = $2`,
		characterID, version)
	return scanVersion(row)
}

func (r *PostgresCharacterRepository) ListVersions(ctx context.Context, characterID domain.CharacterID) ([]*domain.CharacterVersion, error) {
	rows, err := r.q.Query(ctx, `
		SELECT id, character_id, version, appearance, personality, voice_profile, wardrobe
		FROM characters.character_versions WHERE character_id = $1 ORDER BY version`, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.CharacterVersion
	for rows.Next() {
		v, err := scanVersion(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func scanVersion(row rowScanner) (*domain.CharacterVersion, error) {
	var v domain.CharacterVersion
	var appearance, personality, voice, wardrobe []byte
	err := row.Scan(&v.ID, &v.CharacterID, &v.Version, &appearance, &personality, &voice, &wardrobe)
	if postgres.IsNoRows(err) {
		return nil, domain.ErrCharacterNotFound
	}
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(appearance, &v.Appearance)
	_ = json.Unmarshal(personality, &v.Personality)
	_ = json.Unmarshal(voice, &v.VoiceProfile)
	_ = json.Unmarshal(wardrobe, &v.Wardrobe)
	return &v, nil
}

func (r *PostgresCharacterRepository) SaveRelationship(ctx context.Context, characterID domain.CharacterID, rel *domain.Relationship) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO characters.relationships
			(id, project_id, character_id, other_character_id, relation, description)
		VALUES ($1, (SELECT project_id FROM characters.characters WHERE id = $2), $2, $3, $4, $5)`,
		"rel_"+uuid.NewString(), characterID, rel.TargetCharacterID, rel.RelationshipType, rel.Description)
	return err
}

func (r *PostgresCharacterRepository) ListRelationships(ctx context.Context, characterID domain.CharacterID) ([]domain.Relationship, error) {
	rows, err := r.q.Query(ctx, `
		SELECT other_character_id, relation, description
		FROM characters.relationships WHERE character_id = $1`, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Relationship
	for rows.Next() {
		var rel domain.Relationship
		if err := rows.Scan(&rel.TargetCharacterID, &rel.RelationshipType, &rel.Description); err != nil {
			return nil, err
		}
		out = append(out, rel)
	}
	return out, rows.Err()
}

func (r *PostgresCharacterRepository) SaveWardrobeAssignment(ctx context.Context, wa *domain.WardrobeAssignment) error {
	if wa.ID == "" {
		wa.ID = "wa_" + uuid.NewString()
	}
	_, err := r.q.Exec(ctx, `
		INSERT INTO characters.wardrobe_assignments
			(id, character_id, episode_id, scene_id, items, change_event)
		VALUES ($1,$2,$3,$4,$5,$6)`,
		wa.ID, wa.CharacterID, wa.EpisodeID, wa.SceneID, mustJSON(wa.Items), wa.ChangeEvent)
	return err
}

func (r *PostgresCharacterRepository) ListWardrobeAssignments(ctx context.Context, characterID domain.CharacterID, episodeID string) ([]*domain.WardrobeAssignment, error) {
	rows, err := r.q.Query(ctx, `
		SELECT id, character_id, episode_id, scene_id, items, change_event
		FROM characters.wardrobe_assignments
		WHERE character_id = $1 AND ($2 = '' OR episode_id = $2)`,
		characterID, episodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.WardrobeAssignment
	for rows.Next() {
		var wa domain.WardrobeAssignment
		var items []byte
		if err := rows.Scan(&wa.ID, &wa.CharacterID, &wa.EpisodeID, &wa.SceneID, &items, &wa.ChangeEvent); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(items, &wa.Items)
		out = append(out, &wa)
	}
	return out, rows.Err()
}
