package domain

import "errors"

var ErrResetTokenInvalid = errors.New("reset token is invalid, expired, or already used")
