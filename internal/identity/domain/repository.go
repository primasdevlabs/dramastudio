package domain

import (
	"context"
	"time"
)

type UserRepository interface {
	SaveUser(ctx context.Context, u *User) error
	FindUserByID(ctx context.Context, id string) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	ListUsers(ctx context.Context, orgID string) ([]*User, error)
	SaveOrganization(ctx context.Context, o *Organization) error
	FindOrganizationByID(ctx context.Context, id string) (*Organization, error)
	SaveMembership(ctx context.Context, m *Membership) error
	// FindMembershipRole returns the user's role inside the org, or "" when
	// the user is not a member.
	FindMembershipRole(ctx context.Context, userID, orgID string) (string, error)
	SaveCredential(ctx context.Context, c *APICredential) error
	FindCredentialByDigest(ctx context.Context, digest string) (*APICredential, error)
	ListCredentials(ctx context.Context, orgID string) ([]*APICredential, error)
	RevokeCredential(ctx context.Context, id string) error
	// Password-reset tokens are stored as sha256 digests, single-use,
	// time-bounded. FindUserByResetToken only resolves live tokens.
	SaveResetToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	FindUserByResetToken(ctx context.Context, tokenHash string) (*User, error)
	ConsumeResetToken(ctx context.Context, tokenHash string) error
}
