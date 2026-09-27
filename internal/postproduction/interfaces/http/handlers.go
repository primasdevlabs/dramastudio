package http

import "net/http"

type PostproductionHandler struct{}

func NewPostproductionHandler() *PostproductionHandler {
	return &PostproductionHandler{}
}

func (h *PostproductionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
