package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrOrgNotFound        = errors.New("organization not found")
	ErrCredentialNotFound = errors.New("credential not found")
	ErrCredentialRevoked  = errors.New("credential revoked")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailTaken         = errors.New("email already registered")
)
