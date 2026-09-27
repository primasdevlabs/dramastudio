package domain

import "errors"

var (
	ErrModelNotFound          = errors.New("model not found in registry")
	ErrCapabilityNotSupported = errors.New("capability not supported by model")
	ErrPolicyViolation        = errors.New("model policy violation")
	ErrProviderUnavailable    = errors.New("provider unavailable")
)
