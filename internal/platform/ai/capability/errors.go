package capability

import "time"

// ProviderError is the normalized failure surface every adapter must
// produce (Backend.md §84). Workflows decide retry policy from Kind.
type ProviderError struct {
	Kind       ProviderErrorKind
	Provider   string
	Message    string
	Retryable  bool
	RetryAfter time.Duration
}

func (e *ProviderError) Error() string {
	return e.Provider + ": " + e.Kind.String() + ": " + e.Message
}

type ProviderErrorKind string

const (
	ErrRateLimited    ProviderErrorKind = "ProviderRateLimited"
	ErrUnavailable    ProviderErrorKind = "ProviderUnavailable"
	ErrTimeout        ProviderErrorKind = "ProviderTimeout"
	ErrRejected       ProviderErrorKind = "ProviderRejected"
	ErrInvalidRequest ProviderErrorKind = "ProviderInvalidRequest"
	ErrUnknown        ProviderErrorKind = "ProviderError"
)

func (k ProviderErrorKind) String() string { return string(k) }
