package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"dramastudio/internal/identity/domain"
	"dramastudio/internal/platform/audit"
	"dramastudio/internal/platform/security"
)

type IdentityService struct {
	repo  domain.UserRepository
	audit audit.Logger // may be nil; set via SetAudit
}

func NewIdentityService(repo domain.UserRepository) *IdentityService {
	return &IdentityService{repo: repo}
}

// SetAudit injects the append-only audit logger (§63). Nil-safe recorder.
func (s *IdentityService) SetAudit(l audit.Logger) {
	s.audit = l
}

// Register creates an organization + owner user with a bcrypt-hashed
// password.
func (s *IdentityService) Register(ctx context.Context, orgName, email, name, password string) (*domain.User, *domain.Organization, error) {
	if _, err := s.repo.FindUserByEmail(ctx, email); err == nil {
		return nil, nil, domain.ErrEmailTaken
	}
	hash, err := security.HashPassword(password)
	if err != nil {
		return nil, nil, err
	}
	org := &domain.Organization{ID: "org_" + uuid.NewString(), Name: orgName}
	if err := s.repo.SaveOrganization(ctx, org); err != nil {
		return nil, nil, err
	}
	u := &domain.User{
		ID:           "user_" + uuid.NewString(),
		Email:        email,
		Name:         name,
		Role:         domain.RoleOwner,
		OrgID:        org.ID,
		PasswordHash: hash,
		CreatedAt:    time.Now().UTC(),
	}
	if err := s.repo.SaveUser(ctx, u); err != nil {
		return nil, nil, err
	}
	m := &domain.Membership{UserID: u.ID, OrganizationID: org.ID, Role: domain.RoleOwner}
	if err := s.repo.SaveMembership(ctx, m); err != nil {
		return nil, nil, err
	}
	audit.Record(s.audit, ctx, audit.Entry{
		Actor: u.ID, ActorKind: audit.ActorUser, Action: "user.register",
		EntityType: "user", EntityID: u.ID, Detail: map[string]interface{}{"org_id": org.ID},
	})
	return u, org, nil
}

// Authenticate verifies email+password and returns the user on success.
func (s *IdentityService) Authenticate(ctx context.Context, email, password string) (*domain.User, error) {
	u, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}
	if !security.CheckPassword(password, u.PasswordHash) {
		return nil, domain.ErrInvalidCredentials
	}
	audit.Record(s.audit, ctx, audit.Entry{
		Actor: u.ID, ActorKind: audit.ActorUser, Action: "user.login",
		EntityType: "user", EntityID: u.ID,
	})
	return u, nil
}

// AuthenticateAPIKey resolves an API key to its credential via digest.
func (s *IdentityService) AuthenticateAPIKey(ctx context.Context, plaintextKey string) (*domain.APICredential, error) {
	return s.repo.FindCredentialByDigest(ctx, security.APIKeyDigest(plaintextKey))
}

// CreateAPIKey mints a new credential; plaintext is returned once.
func (s *IdentityService) CreateAPIKey(ctx context.Context, orgID, userID, name string, permissions []string, service bool) (*domain.APICredential, string, error) {
	plaintext, digest, err := security.NewAPIKey()
	if err != nil {
		return nil, "", err
	}
	cred := &domain.APICredential{
		ID:          "cred_" + uuid.NewString(),
		OrgID:       orgID,
		UserID:      userID,
		Name:        name,
		KeyDigest:   digest,
		Permissions: permissions,
		Service:     service,
		CreatedAt:   time.Now().UTC(),
	}
	if err := s.repo.SaveCredential(ctx, cred); err != nil {
		return nil, "", err
	}
	audit.Record(s.audit, ctx, audit.Entry{
		Actor: userID, ActorKind: audit.ActorUser, Action: "apikey.create",
		EntityType: "api_credential", EntityID: cred.ID,
		Detail: map[string]interface{}{"org_id": orgID, "name": name},
	})
	return cred, plaintext, nil
}

func (s *IdentityService) ListCredentials(ctx context.Context, orgID string) ([]*domain.APICredential, error) {
	return s.repo.ListCredentials(ctx, orgID)
}

// RevokeCredential revokes a key owned by orgID; cross-org IDs return
// ErrUserNotFound rather than leaking the credential's existence.
func (s *IdentityService) RevokeCredential(ctx context.Context, orgID, id string) error {
	creds, err := s.repo.ListCredentials(ctx, orgID)
	if err != nil {
		return err
	}
	found := false
	for _, c := range creds {
		if c.ID == id {
			found = true
			break
		}
	}
	if !found {
		return domain.ErrUserNotFound
	}
	if err := s.repo.RevokeCredential(ctx, id); err != nil {
		return err
	}
	audit.Record(s.audit, ctx, audit.Entry{
		ActorKind: audit.ActorUser, Action: "apikey.revoke",
		EntityType: "api_credential", EntityID: id,
	})
	return nil
}

// GetUser returns the user only when it belongs to orgID.
func (s *IdentityService) GetUser(ctx context.Context, orgID, id string) (*domain.User, error) {
	u, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u.OrgID != orgID {
		return nil, domain.ErrUserNotFound
	}
	return u, nil
}

func (s *IdentityService) ListUsers(ctx context.Context, orgID string) ([]*domain.User, error) {
	return s.repo.ListUsers(ctx, orgID)
}

// RequestPasswordReset issues a single-use reset token for an existing
// account. Returns ("", nil) for unknown emails so callers can't enumerate
// accounts; the plaintext token is returned for the transport to deliver.
func (s *IdentityService) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	u, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", nil
	}
	token, err := security.RandomToken()
	if err != nil {
		return "", err
	}
	digest := security.APIKeyDigest(token)
	if err := s.repo.SaveResetToken(ctx, u.ID, digest, time.Now().UTC().Add(time.Hour)); err != nil {
		return "", err
	}
	audit.Record(s.audit, ctx, audit.Entry{
		Actor: u.ID, ActorKind: audit.ActorUser, Action: "user.password_reset.request",
		EntityType: "user", EntityID: u.ID,
	})
	return token, nil
}

// ResetPassword consumes a live reset token and sets a new password.
func (s *IdentityService) ResetPassword(ctx context.Context, token, newPassword string) error {
	if len(newPassword) < 8 {
		return domain.ErrResetTokenInvalid
	}
	u, err := s.repo.FindUserByResetToken(ctx, security.APIKeyDigest(token))
	if err != nil {
		return domain.ErrResetTokenInvalid
	}
	hash, err := security.HashPassword(newPassword)
	if err != nil {
		return err
	}
	u.PasswordHash = hash
	if err := s.repo.SaveUser(ctx, u); err != nil {
		return err
	}
	if err := s.repo.ConsumeResetToken(ctx, security.APIKeyDigest(token)); err != nil {
		return err
	}
	audit.Record(s.audit, ctx, audit.Entry{
		Actor: u.ID, ActorKind: audit.ActorUser, Action: "user.password_reset.complete",
		EntityType: "user", EntityID: u.ID,
	})
	return nil
}

// AuthenticateSSO resolves an externally-verified identity to a local user,
// provisioning a new org + owner account on first sign-in (auto account
// detection). orgName is used only when provisioning.
func (s *IdentityService) AuthenticateSSO(ctx context.Context, email, name, orgName string) (*domain.User, bool, error) {
	if u, err := s.repo.FindUserByEmail(ctx, email); err == nil {
		audit.Record(s.audit, ctx, audit.Entry{
			Actor: u.ID, ActorKind: audit.ActorUser, Action: "user.sso_login",
			EntityType: "user", EntityID: u.ID,
		})
		return u, false, nil
	}
	// First SSO sign-in: provision organization + owner. Password is a random
	// unusable secret — the account authenticates via the provider.
	randPw, err := security.RandomToken()
	if err != nil {
		return nil, false, err
	}
	if orgName == "" {
		orgName = name + "'s Studio"
	}
	u, _, err := s.Register(ctx, orgName, email, name, randPw)
	if err != nil {
		return nil, false, err
	}
	return u, true, nil
}
