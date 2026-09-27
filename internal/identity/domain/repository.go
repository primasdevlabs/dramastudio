package domain

import "context"

type UserRepository interface {
	SaveUser(ctx context.Context, u *User) error
	FindUserByID(ctx context.Context, id string) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	ListUsers(ctx context.Context, orgID string) ([]*User, error)
	SaveOrganization(ctx context.Context, o *Organization) error
	FindOrganizationByID(ctx context.Context, id string) (*Organization, error)
	SaveCredential(ctx context.Context, c *APICredential) error
	FindCredentialByDigest(ctx context.Context, digest string) (*APICredential, error)
	ListCredentials(ctx context.Context, orgID string) ([]*APICredential, error)
	RevokeCredential(ctx context.Context, id string) error
}
