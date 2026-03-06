package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
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

// ListDatasets returns the datasets list handler (placeholder)
func ListDatasets(svc interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	}
}

// ListModels returns the models list handler (placeholder)
func ListModels(svc interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	}
}

// CreateEvaluation returns the create evaluation handler (placeholder)
func CreateEvaluation(svc interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"id":"placeholder"}`))
	}
}

// ListEvaluations returns the evaluations list handler (placeholder)
func ListEvaluations(svc interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	}
}

// GetEvaluation returns the get evaluation handler (placeholder)
func GetEvaluation(svc interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"` + id + `"}`))
	}
}

// StreamEvaluation returns the SSE stream handler (placeholder)
func StreamEvaluation(svc interface{}, hub interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Write([]byte("data: placeholder\n\n"))
	}
}
