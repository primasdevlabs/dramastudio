package persistence

import (
	"context"
	"encoding/json"
	"time"

	"dramastudio/internal/identity/domain"
	"dramastudio/internal/platform/database/postgres"
)

// PostgresUserRepository implements domain.UserRepository over PostgreSQL.
// All queries are parameterized (policy: 04-security).
type PostgresUserRepository struct {
	q postgres.Querier
}

func NewPostgresUserRepository(q postgres.Querier) *PostgresUserRepository {
	return &PostgresUserRepository{q: q}
}

func (r *PostgresUserRepository) SaveUser(ctx context.Context, u *domain.User) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO identity.users (id, org_id, email, display_name, password_hash, created_at)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (id) DO UPDATE SET
			org_id = EXCLUDED.org_id,
			display_name = EXCLUDED.display_name,
			password_hash = EXCLUDED.password_hash`,
		u.ID, u.OrgID, u.Email, u.Name, u.PasswordHash, u.CreatedAt)
	return err
}

func (r *PostgresUserRepository) FindUserByID(ctx context.Context, id string) (*domain.User, error) {
	row := r.q.QueryRow(ctx, `
		SELECT u.id, u.org_id, u.email, u.display_name, u.password_hash, u.created_at,
		       COALESCE(m.role, '')
		FROM identity.users u
		LEFT JOIN identity.memberships m ON m.user_id = u.id AND m.org_id = u.org_id
		WHERE u.id = $1`, id)
	return scanUser(row)
}

func (r *PostgresUserRepository) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.q.QueryRow(ctx, `
		SELECT u.id, u.org_id, u.email, u.display_name, u.password_hash, u.created_at,
		       COALESCE(m.role, '')
		FROM identity.users u
		LEFT JOIN identity.memberships m ON m.user_id = u.id AND m.org_id = u.org_id
		WHERE u.email = $1`, email)
	return scanUser(row)
}

func (r *PostgresUserRepository) ListUsers(ctx context.Context, orgID string) ([]*domain.User, error) {
	rows, err := r.q.Query(ctx, `
		SELECT u.id, u.org_id, u.email, u.display_name, u.password_hash, u.created_at,
		       COALESCE(m.role, '')
		FROM identity.users u
		LEFT JOIN identity.memberships m ON m.user_id = u.id AND m.org_id = u.org_id
		WHERE $1 = '' OR u.org_id = $1
		ORDER BY u.created_at`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func (r *PostgresUserRepository) SaveOrganization(ctx context.Context, o *domain.Organization) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO identity.organizations (id, name) VALUES ($1,$2)
		ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name`, o.ID, o.Name)
	return err
}

func (r *PostgresUserRepository) FindOrganizationByID(ctx context.Context, id string) (*domain.Organization, error) {
	var o domain.Organization
	err := r.q.QueryRow(ctx, `SELECT id, name FROM identity.organizations WHERE id = $1`, id).
		Scan(&o.ID, &o.Name)
	if postgres.IsNoRows(err) {
		return nil, domain.ErrOrgNotFound
	}
	return &o, err
}

// SaveMembership upserts the user↔org role binding (§8).
func (r *PostgresUserRepository) SaveMembership(ctx context.Context, m *domain.Membership) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO identity.memberships (user_id, org_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, org_id) DO UPDATE SET role = EXCLUDED.role`,
		m.UserID, m.OrganizationID, m.Role)
	return err
}

// FindMembershipRole returns the member's role or "" for non-members.
func (r *PostgresUserRepository) FindMembershipRole(ctx context.Context, userID, orgID string) (string, error) {
	var role string
	err := r.q.QueryRow(ctx,
		`SELECT role FROM identity.memberships WHERE user_id = $1 AND org_id = $2`,
		userID, orgID).Scan(&role)
	if postgres.IsNoRows(err) {
		return "", nil
	}
	return role, err
}

func (r *PostgresUserRepository) SaveCredential(ctx context.Context, c *domain.APICredential) error {
	perms, err := json.Marshal(c.Permissions)
	if err != nil {
		return err
	}
	_, err = r.q.Exec(ctx, `
		INSERT INTO identity.api_credentials (id, org_id, user_id, name, key_digest, permissions, service, revoked_at, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name, permissions = EXCLUDED.permissions,
			revoked_at = EXCLUDED.revoked_at`,
		c.ID, c.OrgID, c.UserID, c.Name, c.KeyDigest, perms, c.Service, c.RevokedAt, c.CreatedAt)
	return err
}

func (r *PostgresUserRepository) FindCredentialByDigest(ctx context.Context, digest string) (*domain.APICredential, error) {
	row := r.q.QueryRow(ctx, `
		SELECT id, org_id, user_id, name, key_digest, permissions, service, revoked_at, created_at
		FROM identity.api_credentials WHERE key_digest = $1`, digest)
	c, err := scanCredential(row)
	if err != nil {
		return nil, err
	}
	if c.Revoked() {
		return nil, domain.ErrCredentialRevoked
	}
	return c, nil
}

func (r *PostgresUserRepository) ListCredentials(ctx context.Context, orgID string) ([]*domain.APICredential, error) {
	rows, err := r.q.Query(ctx, `
		SELECT id, org_id, user_id, name, key_digest, permissions, service, revoked_at, created_at
		FROM identity.api_credentials WHERE org_id = $1 ORDER BY created_at`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.APICredential
	for rows.Next() {
		c, err := scanCredential(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *PostgresUserRepository) RevokeCredential(ctx context.Context, id string) error {
	tag, err := r.q.Exec(ctx, `
		UPDATE identity.api_credentials SET revoked_at = $2 WHERE id = $1 AND revoked_at IS NULL`,
		id, time.Now().UTC())
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrCredentialNotFound
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanUser(row rowScanner) (*domain.User, error) {
	var u domain.User
	err := row.Scan(&u.ID, &u.OrgID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt, &u.Role)
	if postgres.IsNoRows(err) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func scanCredential(row rowScanner) (*domain.APICredential, error) {
	var c domain.APICredential
	var perms []byte
	var userID *string
	err := row.Scan(&c.ID, &c.OrgID, &userID, &c.Name, &c.KeyDigest, &perms, &c.Service, &c.RevokedAt, &c.CreatedAt)
	if postgres.IsNoRows(err) {
		return nil, domain.ErrCredentialNotFound
	}
	if err != nil {
		return nil, err
	}
	if userID != nil {
		c.UserID = *userID
	}
	if err := json.Unmarshal(perms, &c.Permissions); err != nil {
		c.Permissions = []string{}
	}
	return &c, nil
}

func (r *PostgresUserRepository) SaveResetToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO identity.password_reset_tokens (token_hash, user_id, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (token_hash) DO NOTHING`, tokenHash, userID, expiresAt)
	return err
}

func (r *PostgresUserRepository) FindUserByResetToken(ctx context.Context, tokenHash string) (*domain.User, error) {
	// Expiry is compared against a Go-bound timestamp rather than now() so the
	// check behaves identically under the sqlite dialect translation.
	row := r.q.QueryRow(ctx, `
		SELECT u.id, u.org_id, u.email, u.display_name, u.password_hash, u.created_at,
		       COALESCE(m.role, '')
		FROM identity.password_reset_tokens t
		JOIN identity.users u ON u.id = t.user_id
		LEFT JOIN identity.memberships m ON m.user_id = u.id AND m.org_id = u.org_id
		WHERE t.token_hash = $1 AND t.used_at IS NULL AND t.expires_at > $2`,
		tokenHash, time.Now().UTC())
	return scanUser(row)
}

func (r *PostgresUserRepository) ConsumeResetToken(ctx context.Context, tokenHash string) error {
	_, err := r.q.Exec(ctx, `
		UPDATE identity.password_reset_tokens SET used_at = now()
		WHERE token_hash = $1 AND used_at IS NULL`, tokenHash)
	return err
}
