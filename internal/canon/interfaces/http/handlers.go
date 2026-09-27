package http

import "net/http"

type CanonHandler struct{}

func NewCanonHandler() *CanonHandler {
	return &CanonHandler{}
}

func (h *CanonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
