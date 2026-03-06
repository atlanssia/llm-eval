package stream

import (
	"context"
	"testing"
	"time"
)

func TestHub_Subscribe_Broadcast(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hub := NewHub(ctx, nil)
	defer hub.Close()

	// Subscribe to evaluation
	ch := hub.Subscribe("eval-123")
	defer hub.Unsubscribe("eval-123", ch)

	// Broadcast event
	event := NewProgressEvent("eval-123", "gpt-4", "mmlu", 1, 100)
	hub.Broadcast("eval-123", event)

	// Receive event
	select {
	case received := <-ch:
		if received.Type != EventTypeProgress {
			t.Errorf("expected type 'progress', got '%s'", received.Type)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout waiting for event")
	}
}

func TestHub_Close(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hub := NewHub(ctx, nil)
	ch := hub.Subscribe("eval-123")

	hub.Close()

	// Channel should be closed
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("expected channel to be closed")
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout waiting for close")
	}
}
