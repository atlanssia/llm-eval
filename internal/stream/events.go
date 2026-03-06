package stream

// EventType represents the type of SSE event
type EventType string

const (
	EventTypeProgress           EventType = "progress"
	EventTypeModelComplete      EventType = "model_complete"
	EventTypeEvaluationComplete EventType = "evaluation_complete"
	EventTypeError              EventType = "error"
)

// Event represents an SSE event
type Event struct {
	Type EventType              `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// NewProgressEvent creates a progress event
func NewProgressEvent(evalID, model, dataset string, current, total int) Event {
	var progress float64
	if total > 0 {
		progress = float64(current) / float64(total) * 100
	}

	return Event{
		Type: EventTypeProgress,
		Data: map[string]interface{}{
			"evaluation_id": evalID,
			"model":         model,
			"dataset":       dataset,
			"current":       current,
			"total":         total,
			"progress":      progress,
		},
	}
}

// NewModelCompleteEvent creates a model complete event
func NewModelCompleteEvent(evalID, model string, metrics map[string]interface{}) Event {
	data := map[string]interface{}{
		"evaluation_id": evalID,
		"model":         model,
	}
	for k, v := range metrics {
		data[k] = v
	}
	return Event{
		Type: EventTypeModelComplete,
		Data: data,
	}
}

// NewEvaluationCompleteEvent creates an evaluation complete event
func NewEvaluationCompleteEvent(evalID string, summary map[string]interface{}) Event {
	data := map[string]interface{}{
		"evaluation_id": evalID,
	}
	for k, v := range summary {
		data[k] = v
	}
	return Event{
		Type: EventTypeEvaluationComplete,
		Data: data,
	}
}

// NewErrorEvent creates an error event
func NewErrorEvent(evalID, message string) Event {
	return Event{
		Type: EventTypeError,
		Data: map[string]interface{}{
			"evaluation_id": evalID,
			"error":         message,
		},
	}
}
