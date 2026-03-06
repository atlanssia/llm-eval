package stream

import (
	"encoding/json"
	"testing"
)

func TestEvent_Marshal(t *testing.T) {
	event := Event{
		Type: "progress",
		Data: map[string]interface{}{
			"current": 10,
			"total":   100,
		},
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal event: %v", err)
	}

	var decoded Event
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		t.Fatalf("failed to unmarshal event: %v", err)
	}

	if decoded.Type != "progress" {
		t.Errorf("expected type 'progress', got '%s'", decoded.Type)
	}
}

func TestProgressEvent(t *testing.T) {
	event := NewProgressEvent("eval-123", "gpt-4", "mmlu_anatomy", 10, 100)

	if event.Type != "progress" {
		t.Errorf("expected type 'progress', got '%s'", event.Type)
	}

	data := event.Data
	if data["evaluation_id"] != "eval-123" {
		t.Errorf("expected evaluation_id 'eval-123', got '%v'", data["evaluation_id"])
	}

	if data["model"] != "gpt-4" {
		t.Errorf("expected model 'gpt-4', got '%v'", data["model"])
	}

	if data["dataset"] != "mmlu_anatomy" {
		t.Errorf("expected dataset 'mmlu_anatomy', got '%v'", data["dataset"])
	}

	if data["current"] != 10 {
		t.Errorf("expected current 10, got '%v'", data["current"])
	}

	if data["total"] != 100 {
		t.Errorf("expected total 100, got '%v'", data["total"])
	}

	expectedProgress := 10.0
	if data["progress"] != expectedProgress {
		t.Errorf("expected progress %f, got '%v'", expectedProgress, data["progress"])
	}
}

func TestProgressEvent_DivisionByZero(t *testing.T) {
	// Test case 1: total is 0 (edge case that could cause panic)
	event := NewProgressEvent("eval-123", "gpt-4", "mmlu_anatomy", 5, 0)

	if event.Type != "progress" {
		t.Errorf("expected type 'progress', got '%s'", event.Type)
	}

	data := event.Data
	if data["current"] != 5 {
		t.Errorf("expected current 5, got '%v'", data["current"])
	}

	if data["total"] != 0 {
		t.Errorf("expected total 0, got '%v'", data["total"])
	}

	// Progress should be 0 or a safe default value when total is 0
	progress, ok := data["progress"].(float64)
	if !ok {
		t.Errorf("expected progress to be float64, got '%T'", data["progress"])
	}

	// Progress should be 0 when total is 0 (not NaN or infinity)
	if progress != 0.0 {
		t.Errorf("expected progress 0.0 when total is 0, got '%f'", progress)
	}
}

func TestProgressEvent_EdgeCases(t *testing.T) {
	// Test case 1: current equals total (100% progress)
	event := NewProgressEvent("eval-123", "gpt-4", "mmlu_anatomy", 100, 100)
	if event.Data["progress"] != 100.0 {
		t.Errorf("expected progress 100.0 when current equals total, got '%v'", event.Data["progress"])
	}

	// Test case 2: current is 0 (0% progress)
	event = NewProgressEvent("eval-123", "gpt-4", "mmlu_anatomy", 0, 100)
	if event.Data["progress"] != 0.0 {
		t.Errorf("expected progress 0.0 when current is 0, got '%v'", event.Data["progress"])
	}

	// Test case 3: partial progress (50%)
	event = NewProgressEvent("eval-123", "gpt-4", "mmlu_anatomy", 50, 100)
	if event.Data["progress"] != 50.0 {
		t.Errorf("expected progress 50.0, got '%v'", event.Data["progress"])
	}
}

func TestNewModelCompleteEvent(t *testing.T) {
	metrics := map[string]interface{}{
		"accuracy": 0.95,
		"score":    0.87,
	}

	event := NewModelCompleteEvent("eval-123", "gpt-4", metrics)

	if event.Type != "model_complete" {
		t.Errorf("expected type 'model_complete', got '%s'", event.Type)
	}

	data := event.Data
	if data["evaluation_id"] != "eval-123" {
		t.Errorf("expected evaluation_id 'eval-123', got '%v'", data["evaluation_id"])
	}

	if data["model"] != "gpt-4" {
		t.Errorf("expected model 'gpt-4', got '%v'", data["model"])
	}

	if data["accuracy"] != 0.95 {
		t.Errorf("expected accuracy 0.95, got '%v'", data["accuracy"])
	}

	if data["score"] != 0.87 {
		t.Errorf("expected score 0.87, got '%v'", data["score"])
	}
}

func TestNewModelCompleteEvent_EmptyMetrics(t *testing.T) {
	event := NewModelCompleteEvent("eval-123", "gpt-4", map[string]interface{}{})

	if event.Type != "model_complete" {
		t.Errorf("expected type 'model_complete', got '%s'", event.Type)
	}

	if len(event.Data) != 2 {
		t.Errorf("expected 2 fields in data (evaluation_id and model), got %d", len(event.Data))
	}
}

func TestNewEvaluationCompleteEvent(t *testing.T) {
	summary := map[string]interface{}{
		"total_models": 3,
		"total_tests":  1000,
		"duration_ms":  5000,
	}

	event := NewEvaluationCompleteEvent("eval-123", summary)

	if event.Type != "evaluation_complete" {
		t.Errorf("expected type 'evaluation_complete', got '%s'", event.Type)
	}

	data := event.Data
	if data["evaluation_id"] != "eval-123" {
		t.Errorf("expected evaluation_id 'eval-123', got '%v'", data["evaluation_id"])
	}

	if data["total_models"] != 3 {
		t.Errorf("expected total_models 3, got '%v'", data["total_models"])
	}

	if data["total_tests"] != 1000 {
		t.Errorf("expected total_tests 1000, got '%v'", data["total_tests"])
	}

	if data["duration_ms"] != 5000 {
		t.Errorf("expected duration_ms 5000, got '%v'", data["duration_ms"])
	}
}

func TestNewEvaluationCompleteEvent_EmptySummary(t *testing.T) {
	event := NewEvaluationCompleteEvent("eval-123", map[string]interface{}{})

	if event.Type != "evaluation_complete" {
		t.Errorf("expected type 'evaluation_complete', got '%s'", event.Type)
	}

	if len(event.Data) != 1 {
		t.Errorf("expected 1 field in data (evaluation_id), got %d", len(event.Data))
	}
}

func TestNewErrorEvent(t *testing.T) {
	event := NewErrorEvent("eval-123", "connection timeout")

	if event.Type != "error" {
		t.Errorf("expected type 'error', got '%s'", event.Type)
	}

	data := event.Data
	if data["evaluation_id"] != "eval-123" {
		t.Errorf("expected evaluation_id 'eval-123', got '%v'", data["evaluation_id"])
	}

	if data["error"] != "connection timeout" {
		t.Errorf("expected error message 'connection timeout', got '%v'", data["error"])
	}
}

func TestNewErrorEvent_EmptyMessage(t *testing.T) {
	event := NewErrorEvent("eval-123", "")

	if event.Type != "error" {
		t.Errorf("expected type 'error', got '%s'", event.Type)
	}

	if event.Data["error"] != "" {
		t.Errorf("expected empty error message, got '%v'", event.Data["error"])
	}
}
