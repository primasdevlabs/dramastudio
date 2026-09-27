package http

import (
	"errors"
	"net/http"
)

// DomainError is the common shape modules use for business failures;
// WriteErrorFrom maps them onto the §57 error envelope.
type DomainError struct {
	Code    string
	Message string
	HTTP    int
	Details interface{}
}

func (e *DomainError) Error() string { return e.Message }

func NewDomainError(httpStatus int, code, message string) *DomainError {
	return &DomainError{Code: code, Message: message, HTTP: httpStatus}
}

// WriteErrorFrom maps an error onto the standard error response.
// Unknown errors become 500s without leaking internals (§57).
func WriteErrorFrom(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	var de *DomainError
	if errors.As(err, &de) {
		WriteError(w, de.HTTP, de.Code, de.Message, RequestIDFrom(r), de.Details)
		return
	}
	WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", RequestIDFrom(r), nil)
}

// Convenience constructors shared by handlers.
func ErrNotFound(code, msg string) *DomainError {
	return NewDomainError(http.StatusNotFound, code, msg)
}

func ErrConflict(code, msg string) *DomainError {
	return NewDomainError(http.StatusConflict, code, msg)
}

func ErrBadRequest(code, msg string) *DomainError {
	return NewDomainError(http.StatusBadRequest, code, msg)
}
