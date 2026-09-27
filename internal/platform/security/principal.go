package security

import "context"

// Principal is the authenticated caller attached to a request context.
type Principal struct {
	UserID      string
	OrgID       string
	Roles       []string
	Permissions []string
	Service     bool
}

type principalKey struct{}

// WithPrincipal attaches a principal to the context.
func WithPrincipal(ctx context.Context, p Principal) context.Context {
	return context.WithValue(ctx, principalKey{}, p)
}

// PrincipalFrom returns the authenticated principal, if any.
func PrincipalFrom(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(principalKey{}).(Principal)
	return p, ok
}

// HasPermission reports whether the principal holds a permission. The
// wildcard "*" grants everything (service identities, org owners).
func (p Principal) HasPermission(perm string) bool {
	for _, granted := range p.Permissions {
		if granted == perm || granted == "*" {
			return true
		}
	}
	return false
}
