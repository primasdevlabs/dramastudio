package http

import "net/http"

type WorldHandler struct{}

func NewWorldHandler() *WorldHandler {
	return &WorldHandler{}
}

func (h *WorldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
