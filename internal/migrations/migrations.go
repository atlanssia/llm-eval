package migrations

import (
	"database/sql"
	"embed"
	"fmt"
)

//go:embed *.sql
var migrationFS embed.FS

// RunMigrations runs all database migrations
func RunMigrations(db *sql.DB) error {
	// Read migration file
	content, err := migrationFS.ReadFile("001_init.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute migration
	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	return nil
}
