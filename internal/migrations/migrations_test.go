package migrations

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestRunMigrations(t *testing.T) {
	// Create in-memory database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := RunMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Verify tables exist
	tables := []string{"evaluations", "model_results"}
	for _, table := range tables {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&count)
		if err != nil {
			t.Fatalf("failed to check table %s: %v", table, err)
		}
		if count != 1 {
			t.Errorf("table %s not created", table)
		}
	}

	// Verify indexes exist
	indexes := []string{"idx_model_results_evaluation_id", "idx_evaluations_status", "idx_evaluations_created_at"}
	for _, index := range indexes {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name=?", index).Scan(&count)
		if err != nil {
			t.Fatalf("failed to check index %s: %v", index, err)
		}
		if count != 1 {
			t.Errorf("index %s not created", index)
		}
	}
}
