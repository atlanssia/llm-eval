package model

import (
	"errors"
	"time"
)

// Status represents the evaluation status
type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusCanceled  Status = "canceled"
)

// Evaluation represents an evaluation run
type Evaluation struct {
	ID             string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Status         Status
	Config         EvalConfig
	TotalCases     int
	CompletedCases int
	Results        []ModelResult
	Error          string
}

// ModelResult represents results for a single model
type ModelResult struct {
	ModelName    string
	Predictions  []string
	References   []string
	Latencies    []float64 // seconds
	TokensPerSec []float64
	Metrics      Metrics
	ErrorCount   int
}

// Metrics represents evaluation metrics
type Metrics struct {
	Accuracy           float64
	F1                 float64
	BLEU               float64
	ROUGE_L            float64
	AvgLatency         float64
	AvgTokensPerSecond float64
}

// EvalConfig represents evaluation configuration
type EvalConfig struct {
	Models      []string
	Datasets    []string
	SampleSize  int
	MaxWorkers  int
	Ephemeral   bool // If true, don't persist results
}

// Transition transitions the evaluation status
func (e *Evaluation) Transition(event string) error {
	validTransitions := map[Status]map[string]Status{
		StatusPending: {
			"start":  StatusRunning,
			"cancel": StatusCanceled,
		},
		StatusRunning: {
			"complete": StatusCompleted,
			"fail":     StatusFailed,
			"cancel":   StatusCanceled,
		},
	}

	if nextStates, ok := validTransitions[e.Status]; ok {
		if next, ok := nextStates[event]; ok {
			e.Status = next
			e.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.New("invalid status transition")
}

// Progress returns the completion percentage (0-100)
func (e *Evaluation) Progress() float64 {
	if e.TotalCases == 0 {
		return 0
	}
	return float64(e.CompletedCases) / float64(e.TotalCases) * 100
}
