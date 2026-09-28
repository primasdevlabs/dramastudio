package http

import (
	"context"
	"net/http"
	"strings"

	"dramastudio/internal/platform/security"
)

// ProjectResolver returns the org that owns a project, or an error when the
// project does not exist. Implemented by the composition root over the
// projects service so this package stays module-agnostic.
type ProjectResolver func(ctx context.Context, projectID string) (orgID string, err error)

// ProjectGuard enforces organization ownership on /v1/projects/{id}/...
// routes. It runs after Auth: the URL path is inspected directly (routing has
// not happened yet at middleware time), the owning org is resolved once, and
// the request proceeds only when the principal belongs to that org.
//
// Mismatches and missing projects both answer 404 so existence is not
// leaked across tenants. Projects with org_id = '' (pre-tenancy dev rows)
// remain readable — every project written through the API carries an org.
func ProjectGuard(resolve ProjectResolver) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			projectID := projectIDFromPath(r.URL.Path)
			if projectID == "" {
				next.ServeHTTP(w, r)
				return
			}
			p, ok := security.PrincipalFrom(r.Context())
			if !ok {
				WriteError(w, http.StatusUnauthorized, "UNAUTHENTICATED", "Authentication required", RequestIDFrom(r), nil)
				return
			}
			orgID, err := resolve(r.Context(), projectID)
			if err != nil {
				WriteError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project not found", RequestIDFrom(r), nil)
				return
			}
			if orgID != "" && orgID != p.OrgID {
				WriteError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project not found", RequestIDFrom(r), nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// projectIDFromPath extracts the id segment from /v1/projects/{id}[...]
// paths. Returns "" for /v1/projects (list/create) and non-project routes.
func projectIDFromPath(path string) string {
	const prefix = "/v1/projects/"
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	rest := path[len(prefix):]
	if i := strings.IndexByte(rest, '/'); i >= 0 {
		rest = rest[:i]
	}
	return rest
}
