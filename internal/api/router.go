package api

import (
	"log/slog"

	"github.com/atlanssia/llm-eval/internal/api/handler"
	"github.com/atlanssia/llm-eval/internal/api/middleware"
	"github.com/atlanssia/llm-eval/internal/model"
	"github.com/atlanssia/llm-eval/internal/service"
	"github.com/atlanssia/llm-eval/internal/stream"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// NewRouter creates the HTTP router
func NewRouter(
	evalSvc *service.EvaluationService,
	datasetSvc *service.DatasetService,
	modelSvc *service.ModelService,
	streamHub *stream.Hub,
	cfg *model.Config,
	logger *slog.Logger,
) *chi.Mux {
	r := chi.NewRouter()

	// CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(corsMiddleware.Handler)

	// Logging middleware
	r.Use(middleware.Logger(logger))

	// Recovery middleware
	r.Use(middleware.Recover(logger))

	// Optional auth middleware
	if cfg.Auth.Enabled {
		r.Use(middleware.Auth(cfg.Auth.Password))
	}

	// Health check
	r.Get("/health", handler.Health)

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/datasets", handler.ListDatasets(datasetSvc))
		r.Get("/models", handler.ListModels(modelSvc))

		r.Post("/evaluations", handler.CreateEvaluation(evalSvc))
		r.Get("/evaluations", handler.ListEvaluations(evalSvc))
		r.Get("/evaluations/{id}", handler.GetEvaluation(evalSvc))
		r.Get("/evaluations/{id}/stream", handler.StreamEvaluation(evalSvc, streamHub))
	})

	return r
}
