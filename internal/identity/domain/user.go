package domain

import "time"

// Roles defined by the spec (Backend.md §8). Authorization is
// permission-based; role names are a convenience grouping.
const (
	RoleOwner    = "owner"
	RoleProducer = "producer"
	RoleDirector = "director"
	RoleWriter   = "writer"
	RoleEditor   = "editor"
	RoleReviewer = "reviewer"
	RoleOperator = "operator"
	RoleViewer   = "viewer"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	OrgID        string    `json:"org_id"`
	PasswordHash string    `json:"-"` // never serialized
	CreatedAt    time.Time `json:"created_at"`
}
