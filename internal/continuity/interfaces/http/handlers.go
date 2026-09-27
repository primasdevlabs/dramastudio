package http

import "net/http"

type ContinuityHandler struct{}

func NewContinuityHandler() *ContinuityHandler {
	return &ContinuityHandler{}
}

func (h *ContinuityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
