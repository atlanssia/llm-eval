package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/atlanssia/llm-eval/internal/model"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	if err := RunMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}
	return db
}

func TestEvaluationRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewEvaluation(db, nil)

	eval := &model.Evaluation{
		ID:        "test-123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    model.StatusPending,
		Config: model.EvalConfig{
			Models:     []string{"gpt-4"},
			Datasets:   []string{"mmlu_anatomy"},
			SampleSize: 100,
		},
	}

	err := repo.Create(eval)
	if err != nil {
		t.Fatalf("failed to create evaluation: %v", err)
	}

	// Verify retrieval
	fetched, err := repo.GetByID("test-123")
	if err != nil {
		t.Fatalf("failed to get evaluation: %v", err)
	}

	if fetched.ID != "test-123" {
		t.Errorf("expected ID 'test-123', got '%s'", fetched.ID)
	}
	if fetched.Status != model.StatusPending {
		t.Errorf("expected status 'pending', got '%s'", fetched.Status)
	}
}

func TestEvaluationRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewEvaluation(db, nil)

	eval := &model.Evaluation{
		ID:        "test-123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    model.StatusPending,
		Config:    model.EvalConfig{},
	}

	if err := repo.Create(eval); err != nil {
		t.Fatalf("failed to create evaluation: %v", err)
	}

	// Update status
	eval.Status = model.StatusRunning
	eval.CompletedCases = 10
	eval.TotalCases = 100

	if err := repo.Update(eval); err != nil {
		t.Fatalf("failed to update evaluation: %v", err)
	}

	// Verify update
	fetched, _ := repo.GetByID("test-123")
	if fetched.Status != model.StatusRunning {
		t.Errorf("expected status 'running', got '%s'", fetched.Status)
	}
	if fetched.CompletedCases != 10 {
		t.Errorf("expected completed_cases 10, got %d", fetched.CompletedCases)
	}
}

func TestEvaluationRepository_List(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewEvaluation(db, nil)

	// Create multiple evaluations
	for i := 1; i <= 3; i++ {
		eval := &model.Evaluation{
			ID:        "test-" + string(rune('0'+i)),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Status:    model.StatusPending,
			Config:    model.EvalConfig{},
		}
		if err := repo.Create(eval); err != nil {
			t.Fatalf("failed to create evaluation: %v", err)
		}
	}

	// Test listing with limit
	list, err := repo.List(10, 0)
	if err != nil {
		t.Fatalf("failed to list evaluations: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("expected 3 evaluations, got %d", len(list))
	}
}
