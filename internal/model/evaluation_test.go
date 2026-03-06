package model

import (
	"testing"
)

func TestEvaluationStatus_Transitions(t *testing.T) {
	tests := []struct {
		name     string
		current  Status
		event    string
		expected Status
		valid    bool
	}{
		{
			name:     "pending to running",
			current:  StatusPending,
			event:    "start",
			expected: StatusRunning,
			valid:    true,
		},
		{
			name:     "running to completed",
			current:  StatusRunning,
			event:    "complete",
			expected: StatusCompleted,
			valid:    true,
		},
		{
			name:     "running to failed",
			current:  StatusRunning,
			event:    "fail",
			expected: StatusFailed,
			valid:    true,
		},
		{
			name:     "completed to running is invalid",
			current:  StatusCompleted,
			event:    "start",
			expected: StatusCompleted,
			valid:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eval := &Evaluation{Status: tt.current}
			err := eval.Transition(tt.event)

			if tt.valid && err != nil {
				t.Errorf("expected valid transition, got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Error("expected invalid transition, got nil error")
			}
			if eval.Status != tt.expected {
				t.Errorf("expected status %s, got %s", tt.expected, eval.Status)
			}
		})
	}
}

func TestEvaluation_CalculateProgress(t *testing.T) {
	eval := &Evaluation{
		TotalCases:     100,
		CompletedCases: 45,
	}

	progress := eval.Progress()
	if progress != 45.0 {
		t.Errorf("expected progress 45.0, got %f", progress)
	}
}
