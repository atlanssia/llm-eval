package handler

import (
	"encoding/json"
	"net/http"

	"github.com/atlanssia/llm-eval/internal/service"
)

// ListDatasets lists available datasets
func ListDatasets(datasetSvc *service.DatasetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]string{
			{"name": "mmlu_anatomy", "description": "MMLU Anatomy"},
			{"name": "cmmlu_anatomy", "description": "CMMLU Anatomy"},
		})
	}
}
