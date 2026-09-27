package http

import "net/http"

type PublishingHandler struct{}

func NewPublishingHandler() *PublishingHandler {
	return &PublishingHandler{}
}

func (h *PublishingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
