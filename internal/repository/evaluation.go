package repository

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/atlanssia/llm-eval/internal/model"
)

// EvaluationRepository handles evaluation persistence
type EvaluationRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewEvaluation creates a new evaluation repository
func NewEvaluation(db *sql.DB, logger *slog.Logger) *EvaluationRepository {
	return &EvaluationRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new evaluation
func (r *EvaluationRepository) Create(eval *model.Evaluation) error {
	configJSON, err := json.Marshal(eval.Config)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO evaluations (id, created_at, updated_at, status, config, total_cases, completed_cases, error)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.Exec(query,
		eval.ID,
		eval.CreatedAt.Format(time.RFC3339),
		eval.UpdatedAt.Format(time.RFC3339),
		string(eval.Status),
		string(configJSON),
		eval.TotalCases,
		eval.CompletedCases,
		eval.Error,
	)

	return err
}

// GetByID retrieves an evaluation by ID
func (r *EvaluationRepository) GetByID(id string) (*model.Evaluation, error) {
	query := `
		SELECT id, created_at, updated_at, status, config, total_cases, completed_cases, error
		FROM evaluations
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)

	var eval model.Evaluation
	var configJSON string
	var createdAtStr, updatedAtStr string

	err := row.Scan(
		&eval.ID,
		&createdAtStr,
		&updatedAtStr,
		&eval.Status,
		&configJSON,
		&eval.TotalCases,
		&eval.CompletedCases,
		&eval.Error,
	)

	if err != nil {
		return nil, err
	}

	eval.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	eval.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)

	if err := json.Unmarshal([]byte(configJSON), &eval.Config); err != nil {
		return nil, err
	}

	return &eval, nil
}

// Update updates an evaluation
func (r *EvaluationRepository) Update(eval *model.Evaluation) error {
	query := `
		UPDATE evaluations
		SET updated_at = ?, status = ?, total_cases = ?, completed_cases = ?, error = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		eval.UpdatedAt.Format(time.RFC3339),
		string(eval.Status),
		eval.TotalCases,
		eval.CompletedCases,
		eval.Error,
		eval.ID,
	)

	return err
}

// List retrieves all evaluations
func (r *EvaluationRepository) List(limit, offset int) ([]*model.Evaluation, error) {
	query := `
		SELECT id, created_at, updated_at, status, config, total_cases, completed_cases, error
		FROM evaluations
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var evaluations []*model.Evaluation
	for rows.Next() {
		var eval model.Evaluation
		var configJSON string
		var createdAtStr, updatedAtStr string

		err := rows.Scan(
			&eval.ID,
			&createdAtStr,
			&updatedAtStr,
			&eval.Status,
			&configJSON,
			&eval.TotalCases,
			&eval.CompletedCases,
			&eval.Error,
		)

		if err != nil {
			return nil, err
		}

		eval.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		eval.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)

		if err := json.Unmarshal([]byte(configJSON), &eval.Config); err != nil {
			return nil, err
		}

		evaluations = append(evaluations, &eval)
	}

	return evaluations, rows.Err()
}
