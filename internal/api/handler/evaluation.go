package handler

import (
	"encoding/json"
	"net/http"

	"github.com/atlanssia/llm-eval/internal/service"
	"github.com/atlanssia/llm-eval/internal/stream"
)

// CreateEvaluation creates a new evaluation
func CreateEvaluation(evalSvc *service.EvaluationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "not implemented"})
	}
}

// ListEvaluations lists all evaluations
func ListEvaluations(evalSvc *service.EvaluationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{})
	}
}

// GetEvaluation gets an evaluation by ID
func GetEvaluation(evalSvc *service.EvaluationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "not implemented"})
	}
}

// StreamEvaluation streams evaluation progress via SSE
func StreamEvaluation(evalSvc *service.EvaluationService, hub *stream.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "SSE not supported", http.StatusInternalServerError)
			return
		}

		// TODO: Implement SSE streaming
		flusher.Flush()
	}
}
