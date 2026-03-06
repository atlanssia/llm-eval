package stream

import (
	"context"
	"log/slog"
	"sync"
)

// Hub manages SSE client subscriptions and broadcasts
type Hub struct {
	ctx    context.Context
	logger *slog.Logger
	mu     sync.RWMutex

	// map[evaluation_id][]chan Event
	subs map[string][]chan Event
}

// NewHub creates a new event hub
func NewHub(ctx context.Context, logger *slog.Logger) *Hub {
	return &Hub{
		ctx:    ctx,
		logger: logger,
		subs:   make(map[string][]chan Event),
	}
}

// Subscribe subscribes to events for an evaluation
func (h *Hub) Subscribe(evalID string) chan Event {
	h.mu.Lock()
	defer h.mu.Unlock()

	ch := make(chan Event, 100) // Buffered channel
	h.subs[evalID] = append(h.subs[evalID], ch)

	return ch
}

// Unsubscribe removes a subscription
func (h *Hub) Unsubscribe(evalID string, ch chan Event) {
	h.mu.Lock()
	defer h.mu.Unlock()

	subs := h.subs[evalID]
	for i, s := range subs {
		if s == ch {
			// Remove subscription
			h.subs[evalID] = append(subs[:i], subs[i+1:]...)
			close(ch)
			break
		}
	}

	// Clean up if no more subscribers
	if len(h.subs[evalID]) == 0 {
		delete(h.subs, evalID)
	}
}

// Broadcast sends an event to all subscribers of an evaluation
func (h *Hub) Broadcast(evalID string, event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	subs, ok := h.subs[evalID]
	if !ok {
		return // No subscribers
	}

	for _, ch := range subs {
		select {
		case ch <- event:
			// Sent successfully
		default:
			// Channel full, skip
			if h.logger != nil {
				h.logger.Warn("SSE channel full, dropping event",
					"evaluation_id", evalID,
					"event_type", event.Type,
				)
			}
		}
	}
}

// Close closes all subscriptions and stops the hub
func (h *Hub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for evalID, subs := range h.subs {
		for _, ch := range subs {
			close(ch)
		}
		delete(h.subs, evalID)
	}
}
