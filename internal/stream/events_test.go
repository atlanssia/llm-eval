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
}
