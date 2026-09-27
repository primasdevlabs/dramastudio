package http

import "net/http"

type CharactersHandler struct{}

func NewCharactersHandler() *CharactersHandler {
	return &CharactersHandler{}
}

func (h *CharactersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
