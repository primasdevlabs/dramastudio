package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	maxBodyBytes = 8 << 20 // 8 MiB request body ceiling
)

// DecodeJSON reads a JSON body with a size cap and unknown-field rejection.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid JSON payload: "+err.Error(), RequestIDFrom(r), nil)
		return false
	}
	return true
}

// Validator is implemented by request bodies that validate themselves.
type Validator interface {
	Validate() error
}

// FieldError reports a single validation failure.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationError collects field-level failures (422 response).
type ValidationError struct {
	Fields []FieldError
}

func (e *ValidationError) Error() string {
	parts := make([]string, len(e.Fields))
	for i, f := range e.Fields {
		parts[i] = f.Field + ": " + f.Message
	}
	return strings.Join(parts, "; ")
}

func (e *ValidationError) Add(field, msg string) {
	e.Fields = append(e.Fields, FieldError{Field: field, Message: msg})
}

func (e *ValidationError) HasErrors() bool { return len(e.Fields) > 0 }

// DecodeAndValidate decodes then validates a request body.
func DecodeAndValidate(w http.ResponseWriter, r *http.Request, dst Validator) bool {
	if !DecodeJSON(w, r, dst) {
		return false
	}
	if err := dst.Validate(); err != nil {
		var ve *ValidationError
		if errors.As(err, &ve) {
			WriteError(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Request validation failed", RequestIDFrom(r), map[string]interface{}{"fields": ve.Fields})
			return false
		}
		WriteError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error(), RequestIDFrom(r), nil)
		return false
	}
	return true
}

// Pagination captures list query parameters (§56: pagination + filtering).
type Pagination struct {
	Limit  int
	Cursor string
}

// ParsePagination reads ?limit=&cursor= with a bounded default.
func ParsePagination(r *http.Request) Pagination {
	p := Pagination{Cursor: r.URL.Query().Get("cursor")}
	p.Limit = 50
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			p.Limit = n
		}
	}
	if p.Limit < 1 {
		p.Limit = 1
	}
	if p.Limit > 200 {
		p.Limit = 200
	}
	return p
}

// Page is the standard list envelope.
type Page struct {
	Items      interface{} `json:"items"`
	NextCursor string      `json:"next_cursor,omitempty"`
	Total      *int64      `json:"total,omitempty"`
}

// WritePage writes a paginated list response.
func WritePage(w http.ResponseWriter, items interface{}, nextCursor string, total *int64) {
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"items":       items,
		"next_cursor": nextCursor,
		"total":       total,
	})
}

// IdempotencyKey returns the client-supplied idempotency key, if any (§49).
func IdempotencyKey(r *http.Request) string {
	return strings.TrimSpace(r.Header.Get("Idempotency-Key"))
}

// PathValue is a shorthand for r.PathValue with consistent error shape.
func RequirePathValue(w http.ResponseWriter, r *http.Request, name string) (string, bool) {
	v := r.PathValue(name)
	if v == "" {
		WriteError(w, http.StatusBadRequest, "MISSING_PATH_PARAM", fmt.Sprintf("missing path parameter %q", name), RequestIDFrom(r), nil)
		return "", false
	}
	return v, true
}
