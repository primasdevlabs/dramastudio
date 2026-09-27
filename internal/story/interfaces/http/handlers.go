package http

import "net/http"

type StoryHandler struct{}

func NewStoryHandler() *StoryHandler {
	return &StoryHandler{}
}

func (h *StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
