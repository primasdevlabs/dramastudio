package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting DramaStudio API server...")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
}
