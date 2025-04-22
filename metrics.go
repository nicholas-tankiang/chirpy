package main

import (
	"fmt"
	"net/http"
)

// Implements HTTP handlers for displaying and resetting server access metrics

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	count := cfg.fileserverHits.Load()
	htmlContent := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", count)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlContent))
}
