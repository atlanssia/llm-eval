package repository

import (
	"database/sql"

	"github.com/atlanssia/llm-eval/internal/migrations"
)

// RunMigrations runs all database migrations
func RunMigrations(db *sql.DB) error {
	return migrations.RunMigrations(db)
}
