package http

import "net/http"

type AgentsHandler struct{}

func NewAgentsHandler() *AgentsHandler {
	return &AgentsHandler{}
}

func (h *AgentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
