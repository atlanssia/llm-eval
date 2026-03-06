package repository

import (
	"database/sql"
	"log/slog"
)

// ResultRepository handles result persistence
type ResultRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewResult creates a new result repository
func NewResult(db *sql.DB, logger *slog.Logger) *ResultRepository {
	return &ResultRepository{
		db:     db,
		logger: logger,
	}
}
