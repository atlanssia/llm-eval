package model

import (
	"testing"
)

func TestEvaluationCase_HashID(t *testing.T) {
	cases := []EvaluationCase{
		NewEvaluationCase("What is 2+2?", "4", nil, TaskReasoning, nil),
		NewEvaluationCase("What is 2+2?", "4", nil, TaskReasoning, nil),
		NewEvaluationCase("What is 3+3?", "6", nil, TaskReasoning, nil),
	}

	// Same content should produce same ID
	if cases[0].ID != cases[1].ID {
		t.Error("identical cases should have same ID")
	}

	// Different content should produce different ID
	if cases[0].ID == cases[2].ID {
		t.Error("different cases should have different ID")
	}
}

func TestDataset_Filter(t *testing.T) {
	dataset := Dataset{
		Name: "test",
		Cases: []EvaluationCase{
			{ID: "1", Question: "Q1"},
			{ID: "2", Question: "Q2"},
			{ID: "3", Question: "Q3"},
		},
	}

	filtered := dataset.Filter(2)
	if len(filtered.Cases) != 2 {
		t.Errorf("expected 2 cases, got %d", len(filtered.Cases))
	}
}
