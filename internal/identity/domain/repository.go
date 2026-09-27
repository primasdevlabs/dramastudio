package domain

import "context"

type UserRepository interface {
	SaveUser(ctx context.Context, u *User) error
	FindUserByID(ctx context.Context, id string) (*User, error)
	ListUsers(ctx context.Context) ([]*User, error)
}
