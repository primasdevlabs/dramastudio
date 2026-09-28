package domain

// Membership binds a user to an organization with a role (§8). The
// identity.memberships table is the source of truth for roles — the users
// table intentionally carries no role column.
type Membership struct {
	UserID         string `json:"user_id"`
	OrganizationID string `json:"organization_id"`
	Role           string `json:"role"`
}
