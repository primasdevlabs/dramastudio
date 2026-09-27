package domain

import "time"

// APICredential is a service/organization credential. Only the key digest
// is stored — plaintext keys are shown once at creation (§59).
type APICredential struct {
	ID          string     `json:"id"`
	OrgID       string     `json:"org_id"`
	UserID      string     `json:"user_id,omitempty"`
	Name        string     `json:"name"`
	KeyDigest   string     `json:"-"`
	Permissions []string   `json:"permissions"`
	Service     bool       `json:"service"`
	RevokedAt   *time.Time `json:"revoked_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (c *APICredential) Revoked() bool { return c.RevokedAt != nil }
