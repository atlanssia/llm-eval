package service

import (
	"log/slog"

	"github.com/atlanssia/llm-eval/internal/repository"
	"github.com/atlanssia/llm-eval/internal/stream"
)

// EvaluationService handles evaluation business logic
type EvaluationService struct {
	evalRepo  *repository.EvaluationRepository
	resultRepo *repository.ResultRepository
	hub       *stream.Hub
	logger    *slog.Logger
}

// NewEvaluationService creates a new evaluation service
func NewEvaluationService(
	evalRepo *repository.EvaluationRepository,
	resultRepo *repository.ResultRepository,
	hub *stream.Hub,
	logger *slog.Logger,
) *EvaluationService {
	return &EvaluationService{
		evalRepo:  evalRepo,
		resultRepo: resultRepo,
		hub:       hub,
		logger:    logger,
	}
}
