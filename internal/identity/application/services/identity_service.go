package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"dramastudio/internal/identity/domain"
	"dramastudio/internal/platform/security"
)

type IdentityService struct {
	repo domain.UserRepository
}

func NewIdentityService(repo domain.UserRepository) *IdentityService {
	return &IdentityService{repo: repo}
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
	return cred, plaintext, nil
}

func (s *IdentityService) ListCredentials(ctx context.Context, orgID string) ([]*domain.APICredential, error) {
	return s.repo.ListCredentials(ctx, orgID)
}

func (s *IdentityService) RevokeCredential(ctx context.Context, id string) error {
	return s.repo.RevokeCredential(ctx, id)
}

func (s *IdentityService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.FindUserByID(ctx, id)
}

func (s *IdentityService) ListUsers(ctx context.Context, orgID string) ([]*domain.User, error) {
	return s.repo.ListUsers(ctx, orgID)
}
