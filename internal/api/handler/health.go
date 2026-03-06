package handler

import (
	"encoding/json"
	"net/http"
)

var version = "dev"

// Health returns the health check handler
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"version": version,
	})
}
