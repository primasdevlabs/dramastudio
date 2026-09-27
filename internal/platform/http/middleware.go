package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"dramastudio/internal/platform/observability/logging"
	"dramastudio/internal/platform/security"
)

// Middleware wraps an http.Handler.
type Middleware func(http.Handler) http.Handler

// Chain applies middlewares outermost-first.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

type requestIDKey struct{}

// RequestID assigns/propagates X-Request-Id and stores it in context.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-Id")
		if id == "" {
			id = "req_" + uuid.NewString()
		}
		w.Header().Set("X-Request-Id", id)
		ctx := context.WithValue(r.Context(), requestIDKey{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestIDFrom extracts the request id set by RequestID middleware.
func RequestIDFrom(r *http.Request) string {
	if v, ok := r.Context().Value(requestIDKey{}).(string); ok {
		return v
	}
	return ""
}

// Recoverer converts panics into a consistent 500 API error.
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logging.FromContext(r.Context()).Error("http panic", "error", rec, "path", r.URL.Path)
				WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", RequestIDFrom(r), nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// RequestLogger emits one structured line per request.
func RequestLogger(log logging.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rw, r)
			log.Info("http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.status,
				"duration_ms", time.Since(start).Milliseconds(),
				"request_id", RequestIDFrom(r),
			)
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// Flush preserves http.Flusher through the recorder (needed for SSE).
func (s *statusRecorder) Flush() {
	if f, ok := s.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// AuthConfig controls request authentication behavior.
type AuthConfig struct {
	RequireAuth bool
	DevUserID   string
	DevOrgID    string
}

// Authenticator validates bearer tokens into principals.
type Authenticator interface {
	Validate(token string) (security.Principal, error)
}

// Auth authenticates requests. When RequireAuth is false and no token is
// presented, a development principal is injected so local flows still work.
func Auth(cfg AuthConfig, authn Authenticator) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := bearerToken(r)
			if token == "" {
				if cfg.RequireAuth {
					WriteError(w, http.StatusUnauthorized, "UNAUTHENTICATED", "Authentication required", RequestIDFrom(r), nil)
					return
				}
				p := security.Principal{
					UserID:      cfg.DevUserID,
					OrgID:       cfg.DevOrgID,
					Roles:       []string{"owner"},
					Permissions: []string{"*"},
				}
				next.ServeHTTP(w, r.WithContext(security.WithPrincipal(r.Context(), p)))
				return
			}
			p, err := authn.Validate(token)
			if err != nil {
				WriteError(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token", RequestIDFrom(r), nil)
				return
			}
			next.ServeHTTP(w, r.WithContext(security.WithPrincipal(r.Context(), p)))
		})
	}
}

func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if strings.HasPrefix(h, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
	}
	return ""
}

// RequirePermission authorizes a principal permission per route.
func RequirePermission(perm string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p, ok := security.PrincipalFrom(r.Context())
			if !ok || !p.HasPermission(perm) {
				WriteError(w, http.StatusForbidden, "FORBIDDEN", "Missing permission "+perm, RequestIDFrom(r), nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// CORS applies an explicit allowlist; absent list means same-origin only.
func CORS(allowedOrigins []string) Middleware {
	return func(next http.Handler) http.Handler {
		allowed := make(map[string]bool, len(allowedOrigins))
		for _, o := range allowedOrigins {
			allowed[o] = true
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && allowed[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Idempotency-Key, X-Request-Id")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
