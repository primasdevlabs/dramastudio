package http

import "net/http"

type ProductionHandler struct{}

func NewProductionHandler() *ProductionHandler {
	return &ProductionHandler{}
}

func (h *ProductionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
