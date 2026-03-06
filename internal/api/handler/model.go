package handler

import (
	"encoding/json"
	"net/http"

	"github.com/atlanssia/llm-eval/internal/service"
)

// ListModels lists configured models
func ListModels(modelSvc *service.ModelService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]string{
			{"name": "gpt-4", "endpoint": "https://api.openai.com/v1/chat/completions"},
		})
	}
}
