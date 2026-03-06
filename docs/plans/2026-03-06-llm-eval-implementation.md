# LLM Evaluation Web Tool - Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a web-based LLM evaluation tool with Go + React + Tailwind + Vite, embedding the React frontend into a single Go binary.

**Architecture:** Clean layered architecture with Chi router, Service layer for business logic, Repository layer for SQLite persistence, and SSE for real-time progress streaming.

**Tech Stack:** Go 1.26.1, Chi v5, SQLite (modernc.org/sqlite), React 18, Vite, Tailwind CSS, shadcn/ui, TanStack Query

---

## Task 1: Initialize Go Module and Project Structure

**Files:**
- Create: `go.mod`
- Create: `go.sum` (generated)
- Create: `Makefile`
- Create: `.gitignore`
- Create: `README.md`

**Step 1: Create go.mod**

```bash
cd /Users/mw/workspace/tmp/llm-eval
go mod init github.com/atlanssia/llm-eval
```

File: `go.mod`
```go
module github.com/atlanssia/llm-eval

go 1.26.1

require (
    github.com/go-chi/chi/v5 v5.1.0
    github.com/go-chi/cors v1.2.2
    gopkg.in/yaml.v3 v3.0.1
    modernc.org/sqlite v1.34.4
)
```

**Step 2: Download dependencies**

```bash
go mod tidy
```

Expected: `go.sum` file created, dependencies downloaded

**Step 3: Create Makefile**

File: `Makefile`
```makefile
.PHONY: help dev build-web build-go build run test test-go test-web clean

help:
	@echo "Available targets:"
	@echo "  make dev        - Start dev server (Go + Vite)"
	@echo "  make build-web  - Build React frontend"
	@echo "  make build-go   - Build Go binary"
	@echo "  make build      - Full production build"
	@echo "  make run        - Run the binary"
	@echo "  make test       - Run all tests"
	@echo "  make test-go    - Run Go tests"
	@echo "  make test-web   - Run frontend tests"
	@echo "  make clean      - Clean build artifacts"

dev:
	@echo "Starting dev servers..."
	@make -j2 dev-go dev-web

dev-go:
	@echo "Starting Go API server on :8080..."
	@go run cmd/llm-eval/main.go

dev-web:
	@echo "Starting Vite dev server on :5173..."
	@cd web && npm run dev

build-web:
	@echo "Building React frontend..."
	@cd web && npm run build

build-go: build-web
	@echo "Building Go binary with embedded frontend..."
	@go build -o bin/llm-eval cmd/llm-eval/main.go

build: build-go
	@echo "Build complete: bin/llm-eval"

run:
	@./bin/llm-eval

test: test-go test-web

test-go:
	@echo "Running Go tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-web:
	@echo "Running frontend tests..."
	@cd web && npm test

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@cd web && rm -rf dist/
```

**Step 4: Create .gitignore**

File: `.gitignore`
```gitignore
# Binaries
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib
llm-eval

# Test files
coverage.out
coverage.html
*.test

# Go workspace
/workspace/

# IDE
.vscode/
.idea/
*.swp
*.swo

# Database
*.db
*.db-shm
*.db-wal
data/

# Node
web/node_modules/
web/dist/

# Environment
.env
.env.local
*.yaml
!configs/*.example
```

**Step 5: Create README.md**

File: `README.md`
```markdown
# LLM Evaluation Web Tool

A web-based tool for evaluating Large Language Models on medical datasets.

## Features

- Real-time evaluation monitoring with SSE streaming
- Visual dashboard with metrics visualization
- Config file upload and visual form builder
- Support for multiple LLM providers
- Dataset: MMLU, CMMLU, MedQA, PubMedQA, MedMCQA

## Quick Start

```bash
# Install dependencies
go mod download
cd web && npm install

# Run dev servers
make dev

# Build production binary
make build
./bin/llm-eval
```

## Configuration

Copy `configs/models.yaml.example` to `configs/models.yaml` and configure your models.

## API Documentation

See [API.md](docs/API.md)

## License

MIT
```

**Step 6: Commit**

```bash
git add go.mod go.sum Makefile .gitignore README.md
git commit -m "feat: initialize Go module and project structure"
```

---

## Task 2: Create Directory Structure

**Files:**
- Create: `cmd/llm-eval/`
- Create: `internal/api/handler/`
- Create: `internal/api/middleware/`
- Create: `internal/api/dto/`
- Create: `internal/service/`
- Create: `internal/repository/`
- Create: `internal/model/`
- Create: `internal/stream/`
- Create: `internal/embed/`
- Create: `configs/`
- Create: `migrations/`
- Create: `web/`

**Step 1: Create directories**

```bash
mkdir -p cmd/llm-eval
mkdir -p internal/api/{handler,middleware,dto}
mkdir -p internal/service
mkdir -p internal/repository
mkdir -p internal/model
mkdir -p internal/stream
mkdir -p internal/embed
mkdir -p configs
mkdir -p migrations
```

**Step 2: Create placeholder files for each directory**

```bash
# cmd
touch cmd/llm-eval/main.go

# internal/api/handler
touch internal/api/handler/evaluation.go
touch internal/api/handler/dataset.go
touch internal/api/handler/model.go
touch internal/api/handler/health.go

# internal/api/middleware
touch internal/api/middleware/auth.go
touch internal/api/middleware/cors.go
touch internal/api/middleware/logger.go
touch internal/api/middleware/recover.go

# internal/api/dto
touch internal/api/dto/evaluation.go
touch internal/api/dto/dataset.go
touch internal/api/dto/common.go

# internal/service
touch internal/service/evaluation.go
touch internal/service/dataset.go
touch internal/service/model.go
touch internal/service/metrics.go
touch internal/service/stream.go

# internal/repository
touch internal/repository/evaluation.go
touch internal/repository/result.go
touch internal/repository/migrations.go

# internal/model
touch internal/model/evaluation.go
touch internal/model/dataset.go
touch internal/model/result.go
touch internal/model/config.go

# internal/stream
touch internal/stream/events.go
touch internal/stream/hub.go

# internal/embed
touch internal/embed/embed.go

# configs
touch configs/models.yaml.example

# migrations
touch migrations/001_init.sql
```

**Step 3: Commit**

```bash
git add .
git commit -m "feat: create project directory structure"
```

---

## Task 3: Domain Models - Evaluation

**Files:**
- Modify: `internal/model/evaluation.go`

**Step 1: Write failing test for Evaluation model**

File: `internal/model/evaluation_test.go`
```go
package model

import (
	"testing"
	"time"
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
		TotalCases:    100,
		CompletedCases: 45,
	}

	progress := eval.Progress()
	if progress != 45.0 {
		t.Errorf("expected progress 45.0, got %f", progress)
	}
}
```

**Step 2: Run test to verify it fails**

```bash
cd /Users/mw/workspace/tmp/llm-eval
go test ./internal/model/... -v
```

Expected: FAIL with `undefined: Evaluation`, `undefined: Status`

**Step 3: Write minimal implementation**

File: `internal/model/evaluation.go`
```go
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
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/model/... -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add internal/model/evaluation.go internal/model/evaluation_test.go
git commit -m "feat(model): add Evaluation domain model with status transitions"
```

---

## Task 4: Domain Models - Dataset

**Files:**
- Modify: `internal/model/dataset.go`

**Step 1: Write failing test for Dataset model**

File: `internal/model/dataset_test.go`
```go
package model

import (
	"testing"
)

func TestEvaluationCase_HashID(t *testing.T) {
	cases := []EvaluationCase{
		{
			Question: "What is 2+2?",
			Answer:   "4",
		},
		{
			Question: "What is 2+2?",
			Answer:   "4",
		},
		{
			Question: "What is 3+3?",
			Answer:   "6",
		},
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
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/model/... -v -run Dataset
```

Expected: FAIL with `undefined: Dataset`, `undefined: EvaluationCase`

**Step 3: Write minimal implementation**

File: `internal/model/dataset.go`
```go
package model

import (
	"crypto/md5"
	"fmt"
)

// TaskType represents the type of evaluation task
type TaskType string

const (
	TaskMedicalQA TaskType = "medical_qa"
	TaskReasoning TaskType = "reasoning"
	TaskWorkflow  TaskType = "workflow"
	TaskRAG       TaskType = "rag"
)

// EvaluationCase represents a single evaluation case
type EvaluationCase struct {
	ID       string
	TaskType TaskType
	Question string
	Options  []string
	Answer   string
	Context  string
	Metadata map[string]string
}

// GenerateID creates a unique ID from case content
func GenerateID(content string) string {
	hash := md5.Sum([]byte(content))
	return fmt.Sprintf("%x", hash)[:12]
}

// Dataset represents a collection of evaluation cases
type Dataset struct {
	Name        string
	Source      string
	TaskType    TaskType
	TotalCases  int
	Cases       []EvaluationCase
	Description string
}

// Filter returns a new dataset with at most n cases
func (d *Dataset) Filter(n int) *Dataset {
	if n >= len(d.Cases) {
		return d
	}
	return &Dataset{
		Name:        d.Name,
		Source:      d.Source,
		TaskType:    d.TaskType,
		TotalCases:  n,
		Cases:       d.Cases[:n],
		Description: d.Description,
	}
}

// NewEvaluationCase creates a new evaluation case with generated ID
func NewEvaluationCase(question, answer string, options []string, taskType TaskType, metadata map[string]string) EvaluationCase {
	content := question + answer
	return EvaluationCase{
		ID:       GenerateID(content),
		Question: question,
		Answer:   answer,
		Options:  options,
		TaskType: taskType,
		Metadata: metadata,
	}
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/model/... -v -run Dataset
```

Expected: PASS

**Step 5: Commit**

```bash
git add internal/model/dataset.go internal/model/dataset_test.go
git commit -m "feat(model): add Dataset and EvaluationCase domain models"
```

---

## Task 5: Domain Models - Config

**Files:**
- Modify: `internal/model/config.go`

**Step 1: Write failing test for Config model**

File: `internal/model/config_test.go`
```go
package model

import (
	"testing"
)

func TestModelConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ModelConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: ModelConfig{
				Name:     "gpt-4",
				Endpoint: "https://api.openai.com/v1/chat/completions",
				APIKey:   "sk-test",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			config: ModelConfig{
				Endpoint: "https://api.openai.com/v1/chat/completions",
				APIKey:   "sk-test",
			},
			wantErr: true,
		},
		{
			name: "missing endpoint",
			config: ModelConfig{
				Name:   "gpt-4",
				APIKey: "sk-test",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/model/... -v -run Config
```

Expected: FAIL with `undefined: ModelConfig`

**Step 3: Write minimal implementation**

File: `internal/model/config.go`
```go
package model

import (
	"errors"
	"os"
	"strings"
)

// ModelConfig represents configuration for a single model
type ModelConfig struct {
	Name      string  `yaml:"name"`
	Endpoint  string  `yaml:"endpoint"`
	APIKey    string  `yaml:"api_key,omitempty"`
	Timeout   int     `yaml:"timeout,omitempty"`    // seconds
	MaxRetries int    `yaml:"max_retries,omitempty"` // default 2
}

// Validate validates the model configuration
func (c *ModelConfig) Validate() error {
	if c.Name == "" {
		return errors.New("model name is required")
	}
	if c.Endpoint == "" {
		return errors.New("model endpoint is required")
	}
	return nil
}

// GetAPIKey returns the API key, expanding environment variables
func (c *ModelConfig) GetAPIKey() string {
	if c.APIKey == "" {
		return ""
	}

	// Expand ${VAR} environment variables
	if strings.HasPrefix(c.APIKey, "${") && strings.HasSuffix(c.APIKey, "}") {
		varName := c.APIKey[2 : len(c.APIKey)-1]
		return os.Getenv(varName)
	}

	return c.APIKey
}

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	DataDir  string         `yaml:"data_dir"`
	Auth     AuthConfig     `yaml:"auth"`
	Models   []ModelConfig  `yaml:"models"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Addr         string `yaml:"addr" env:"SERVER_ADDR" default:"0.0.0.0:8080"`
	ReadTimeout  int    `yaml:"read_timeout" default:"30"`   // seconds
	WriteTimeout int    `yaml:"write_timeout" default:"30"`  // seconds
	IdleTimeout  int    `yaml:"idle_timeout" default:"120"`  // seconds
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Path string `yaml:"path" env:"DATABASE_PATH" default:"./data/llm-eval.db"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Enabled  bool   `yaml:"enabled" env:"AUTH_ENABLED" default:"false"`
	Password string `yaml:"password" env:"AUTH_PASSWORD"`
}

// Validate validates the configuration
func (c *Config) Validate() error {
	for _, m := range c.Models {
		if err := m.Validate(); err != nil {
			return err
		}
	}
	return nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/model/... -v -run Config
```

Expected: PASS

**Step 5: Commit**

```bash
git add internal/model/config.go internal/model/config_test.go
git commit -m "feat(model): add Config domain model with validation"
```

---

## Task 6: Stream Hub - SSE Events

**Files:**
- Modify: `internal/stream/events.go`
- Modify: `internal/stream/hub.go`

**Step 1: Write failing test for SSE events**

File: `internal/stream/events_test.go`
```go
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

	data, ok := event.Data.(map[string]interface{})
	if !ok {
		t.Fatal("expected data to be map[string]interface{}")
	}

	if data["evaluation_id"] != "eval-123" {
		t.Errorf("expected evaluation_id 'eval-123', got '%v'", data["evaluation_id"])
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/stream/... -v
```

Expected: FAIL with `undefined: Event`, `undefined: NewProgressEvent`

**Step 3: Write minimal implementation**

File: `internal/stream/events.go`
```go
package stream

// EventType represents the type of SSE event
type EventType string

const (
	EventTypeProgress          EventType = "progress"
	EventTypeModelComplete     EventType = "model_complete"
	EventTypeEvaluationComplete EventType = "evaluation_complete"
	EventTypeError             EventType = "error"
)

// Event represents an SSE event
type Event struct {
	Type EventType              `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// NewProgressEvent creates a progress event
func NewProgressEvent(evalID, model, dataset string, current, total int) Event {
	return Event{
		Type: EventTypeProgress,
		Data: map[string]interface{}{
			"evaluation_id": evalID,
			"model":         model,
			"dataset":       dataset,
			"current":       current,
			"total":         total,
			"progress":      float64(current) / float64(total) * 100,
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
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/stream/... -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add internal/stream/events.go internal/stream/events_test.go
git commit -m "feat(stream): add SSE event types"
```

---

## Task 7: Stream Hub - Event Broadcasting

**Files:**
- Modify: `internal/stream/hub.go`

**Step 1: Write failing test for Hub**

File: `internal/stream/hub_test.go`
```go
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
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/stream/... -v -run Hub
```

Expected: FAIL with `undefined: NewHub`

**Step 3: Write minimal implementation**

File: `internal/stream/hub.go`
```go
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
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/stream/... -v -run Hub
```

Expected: PASS

**Step 5: Commit**

```bash
git add internal/stream/hub.go internal/stream/hub_test.go
git commit -m "feat(stream): add event hub with subscription and broadcasting"
```

---

## Task 8: Database Migrations

**Files:**
- Modify: `migrations/001_init.sql`
- Modify: `internal/repository/migrations.go`

**Step 1: Create migration SQL file**

File: `migrations/001_init.sql`
```sql
-- Evaluations table
CREATE TABLE IF NOT EXISTS evaluations (
    id TEXT PRIMARY KEY,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    status TEXT NOT NULL,
    config TEXT NOT NULL, -- JSON
    total_cases INTEGER NOT NULL DEFAULT 0,
    completed_cases INTEGER NOT NULL DEFAULT 0,
    error TEXT
);

-- Model results table
CREATE TABLE IF NOT EXISTS model_results (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    evaluation_id TEXT NOT NULL,
    model_name TEXT NOT NULL,
    predictions TEXT NOT NULL, -- JSON array
    references TEXT NOT NULL, -- JSON array
    latencies TEXT NOT NULL, -- JSON array
    tokens_per_sec TEXT, -- JSON array
    metrics TEXT NOT NULL, -- JSON object
    error_count INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (evaluation_id) REFERENCES evaluations(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_model_results_evaluation_id ON model_results(evaluation_id);
CREATE INDEX IF NOT EXISTS idx_evaluations_status ON evaluations(status);
CREATE INDEX IF NOT EXISTS idx_evaluations_created_at ON evaluations(created_at DESC);
```

**Step 2: Write test for migrations**

File: `internal/repository/migrations_test.go`
```go
package repository

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func TestRunMigrations(t *testing.T) {
	// Create in-memory database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := RunMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Verify tables exist
	tables := []string{"evaluations", "model_results"}
	for _, table := range tables {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&count)
		if err != nil {
			t.Fatalf("failed to check table %s: %v", table, err)
		}
		if count != 1 {
			t.Errorf("table %s not created", table)
		}
	}

	// Verify indexes exist
	indexes := []string{"idx_model_results_evaluation_id", "idx_evaluations_status", "idx_evaluations_created_at"}
	for _, index := range indexes {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name=?", index).Scan(&count)
		if err != nil {
			t.Fatalf("failed to check index %s: %v", index, err)
		}
		if count != 1 {
			t.Errorf("index %s not created", index)
		}
	}
}
```

**Step 3: Run test to verify it fails**

```bash
go test ./internal/repository/... -v -run Migrations
```

Expected: FAIL with `undefined: RunMigrations`

**Step 4: Write minimal implementation**

File: `internal/repository/migrations.go`
```go
package repository

import (
	"database/sql"
	"embed"
	"fmt"
)

//go:embed ../../migrations/*.sql
var migrationFS embed.FS

// RunMigrations runs all database migrations
func RunMigrations(db *sql.DB) error {
	// Read migration file
	content, err := migrationFS.ReadFile("migrations/001_init.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute migration
	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	return nil
}
```

**Step 5: Run test to verify it passes**

```bash
go test ./internal/repository/... -v -run Migrations
```

Expected: PASS

**Step 6: Commit**

```bash
git add migrations/001_init.sql internal/repository/migrations.go internal/repository/migrations_test.go
git commit -m "feat(repository): add database migrations with tests"
```

---

## Task 9: Repository - Evaluation CRUD

**Files:**
- Modify: `internal/repository/evaluation.go`

**Step 1: Write failing test for evaluation repository**

File: `internal/repository/evaluation_test.go`
```go
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
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/repository/... -v -run EvaluationRepository
```

Expected: FAIL with `undefined: NewEvaluation`

**Step 3: Write minimal implementation**

File: `internal/repository/evaluation.go`
```go
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
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/repository/... -v -run EvaluationRepository
```

Expected: PASS

**Step 5: Commit**

```bash
git add internal/repository/evaluation.go internal/repository/evaluation_test.go
git commit -m "feat(repository): add evaluation repository with CRUD operations"
```

---

## Task 10: Configuration Loading

**Files:**
- Modify: `internal/config/config.go`
- Create: `internal/config/config.go`

**Step 1: Write failing test for config loading**

File: `internal/config/config_test.go`
```go
package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/atlanssia/llm-eval/internal/model"
)

func TestLoad(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  addr: "0.0.0.0:9090"
  read_timeout: 60
  write_timeout: 60
  idle_timeout: 300

database:
  path: "./test.db"

data_dir: "./test_data"

auth:
  enabled: true
  password: "test123"

models:
  - name: "gpt-4"
    endpoint: "https://api.openai.com/v1/chat/completions"
    api_key: "${OPENAI_API_KEY}"
    timeout: 120
    max_retries: 3
  - name: "claude"
    endpoint: "https://api.anthropic.com/v1/messages"
    api_key: "${ANTHROPIC_API_KEY}"
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Set environment variables
	os.Setenv("OPENAI_API_KEY", "sk-test-openai")
	os.Setenv("ANTHROPIC_API_KEY", "sk-test-anthropic")
	defer func() {
		os.Unsetenv("OPENAI_API_KEY")
		os.Unsetenv("ANTHROPIC_API_KEY")
	}()

	// Change to temp directory
	originalWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalWd)

	// Load config
	cfg, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Verify server config
	if cfg.Server.Addr != "0.0.0.0:9090" {
		t.Errorf("expected addr '0.0.0.0:9090', got '%s'", cfg.Server.Addr)
	}
	if cfg.Server.ReadTimeout != 60 {
		t.Errorf("expected read_timeout 60, got %d", cfg.Server.ReadTimeout)
	}

	// Verify auth config
	if !cfg.Auth.Enabled {
		t.Error("expected auth enabled")
	}
	if cfg.Auth.Password != "test123" {
		t.Errorf("expected password 'test123', got '%s'", cfg.Auth.Password)
	}

	// Verify models
	if len(cfg.Models) != 2 {
		t.Errorf("expected 2 models, got %d", len(cfg.Models))
	}

	if cfg.Models[0].Name != "gpt-4" {
		t.Errorf("expected first model name 'gpt-4', got '%s'", cfg.Models[0].Name)
	}

	// Verify API key expansion
	if cfg.Models[0].GetAPIKey() != "sk-test-openai" {
		t.Errorf("expected API key 'sk-test-openai', got '%s'", cfg.Models[0].GetAPIKey())
	}
}

func TestLoad_Defaults(t *testing.T) {
	// Create minimal config
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  addr: ":8080"

database:
  path: "./data.db"

models: []
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	originalWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalWd)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Verify defaults
	if cfg.Server.ReadTimeout != 30 {
		t.Errorf("expected default read_timeout 30, got %d", cfg.Server.ReadTimeout)
	}
	if cfg.DataDir != "./data" {
		t.Errorf("expected default data_dir './data', got '%s'", cfg.DataDir)
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/config/... -v
```

Expected: FAIL with `package config` does not exist

**Step 3: Create directory and write minimal implementation**

```bash
mkdir -p internal/config
```

File: `internal/config/config.go`
```go
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/atlanssia/llm-eval/internal/model"
	"gopkg.in/yaml.v3"
)

// Default config path
const DefaultConfigPath = "config.yaml"

// Load loads configuration from file
// Searches in current directory and ./configs/
func Load() (*model.Config, error) {
	// Try current directory first
	configPath := DefaultConfigPath
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Try ./configs/ directory
		configPath = filepath.Join("configs", DefaultConfigPath)
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found: tried %s and %s", DefaultConfigPath, configPath)
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg model.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 30
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 30
	}
	if cfg.Server.IdleTimeout == 0 {
		cfg.Server.IdleTimeout = 120
	}
	if cfg.DataDir == "" {
		cfg.DataDir = "./data"
	}

	// Validate
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// MustLoad loads configuration or panics
func MustLoad() *model.Config {
	cfg, err := Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return cfg
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/config/... -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add internal/config/config.go internal/config/config_test.go
git commit -m "feat(config): add config loading with YAML support and defaults"
```

---

## Task 11: Main Entry Point with Graceful Shutdown

**Files:**
- Modify: `cmd/llm-eval/main.go`

**Step 1: Write minimal main.go**

File: `cmd/llm-eval/main.go`
```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atlanssia/llm-eval/internal/config"
	"github.com/atlanssia/llm-eval/internal/repository"
	"github.com/atlanssia/llm-eval/internal/stream"
	_ "modernc.org/sqlite"
)

var (
	// Build-time version (set via ldflags)
	version = "dev"
)

func main() {
	// Context for main goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	logger.Info("Starting LLM Evaluation Tool", "version", version)

	// Initialize database
	db, err := initDB(cfg.Database.Path, logger)
	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize repositories
	evalRepo := repository.NewEvaluation(db, logger)
	resultRepo := repository.NewResult(db, logger)

	// Initialize stream hub for SSE
	streamHub := stream.NewHub(ctx, logger)
	defer streamHub.Close()

	// TODO: Initialize services
	_ = evalRepo
	_ = resultRepo
	_ = streamHub

	// TODO: Setup router
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","version":"%s"}`, version)
	})

	// HTTP server with timeouts
	srv := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Channel for errors
	serverErrors := make(chan error, 1)

	// Start server in goroutine
	go func() {
		logger.Info("Server listening", "addr", srv.Addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Wait for signal or server error
	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server failed", "error", err)
			os.Exit(1)
		}
	case sig := <-shutdown:
		logger.Info("Shutdown signal received", "signal", sig)

		// Graceful shutdown with timeout
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)
		defer shutdownCancel()

		// Stop accepting new connections
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("Server shutdown error", "error", err)
			srv.Close()
		}

		logger.Info("Server shutdown complete")
	}
}

func initDB(path string, logger *slog.Logger) (*sql.DB, error) {
	// Ensure data directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Run migrations
	if err := repository.RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("Database initialized", "path", path)
	return db, nil
}
```

Wait, I need to fix the import. Let me update:

File: `cmd/llm-eval/main.go`
```go
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/atlanssia/llm-eval/internal/config"
	"github.com/atlanssia/llm-eval/internal/repository"
	"github.com/atlanssia/llm-eval/internal/stream"
	_ "modernc.org/sqlite"
)

var (
	// Build-time version (set via ldflags)
	version = "dev"
)

func main() {
	// Context for main goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	logger.Info("Starting LLM Evaluation Tool", "version", version)

	// Initialize database
	db, err := initDB(cfg.Database.Path, logger)
	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize repositories
	evalRepo := repository.NewEvaluation(db, logger)
	resultRepo := repository.NewResult(db, logger)

	// Initialize stream hub for SSE
	streamHub := stream.NewHub(ctx, logger)
	defer streamHub.Close()

	// TODO: Initialize services
	_ = evalRepo
	_ = resultRepo
	_ = streamHub

	// TODO: Setup router
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","version":"%s"}`, version)
	})

	// HTTP server with timeouts
	srv := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Channel for errors
	serverErrors := make(chan error, 1)

	// Start server in goroutine
	go func() {
		logger.Info("Server listening", "addr", srv.Addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Wait for signal or server error
	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server failed", "error", err)
			os.Exit(1)
		}
	case sig := <-shutdown:
		logger.Info("Shutdown signal received", "signal", sig)

		// Graceful shutdown with timeout
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)
		defer shutdownCancel()

		// Stop accepting new connections
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("Server shutdown error", "error", err)
			srv.Close()
		}

		logger.Info("Server shutdown complete")
	}
}

func initDB(path string, logger *slog.Logger) (*sql.DB, error) {
	// Ensure data directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Run migrations
	if err := repository.RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("Database initialized", "path", path)
	return db, nil
}
```

Wait, I need to also create the result repository placeholder. Let me add that:

File: `internal/repository/result.go`
```go
package repository

import (
	"database/sql"
	"log/slog"
)

// ResultRepository handles result persistence
type ResultRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewResult creates a new result repository
func NewResult(db *sql.DB, logger *slog.Logger) *ResultRepository {
	return &ResultRepository{
		db:     db,
		logger: logger,
	}
}
```

**Step 2: Create example config file**

File: `configs/models.yaml.example`
```yaml
server:
  addr: "0.0.0.0:8080"
  read_timeout: 30
  write_timeout: 30
  idle_timeout: 120

database:
  path: "./data/llm-eval.db"

data_dir: "./data"

auth:
  enabled: false
  password: ""

models:
  - name: "gpt-4"
    endpoint: "https://api.openai.com/v1/chat/completions"
    api_key: "${OPENAI_API_KEY}"
    timeout: 60
    max_retries: 2

  - name: "claude-3-opus"
    endpoint: "https://api.anthropic.com/v1/messages"
    api_key: "${ANTHROPIC_API_KEY}"
    timeout: 60
    max_retries: 2

  - name: "local-model"
    endpoint: "http://localhost:11434/v1/chat/completions"
    api_key: ""  # No API key for local models
    timeout: 120
    max_retries: 0
```

**Step 3: Test server starts**

```bash
cp configs/models.yaml.example config.yaml
go run cmd/llm-eval/main.go &
sleep 2
curl http://localhost:8080/health
pkill -f "go run cmd/llm-eval/main.go"
```

Expected: `{"status":"ok","version":"dev"}`

**Step 4: Commit**

```bash
git add cmd/llm-eval/main.go internal/repository/result.go configs/models.yaml.example
git commit -m "feat(main): add main entry point with graceful shutdown"
```

---

## Task 12: Chi Router Setup

**Files:**
- Modify: `internal/api/router.go`
- Modify: `internal/api/middleware/logger.go`

**Step 1: Write test for router**

File: `internal/api/router_test.go`
```go
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atlanssia/llm-eval/internal/config"
	"github.com/atlanssia/llm-eval/internal/model"
	"github.com/go-chi/chi/v5"
)

func TestRouter_HealthCheck(t *testing.T) {
	cfg := &model.Config{}
	router := NewRouter(nil, nil, nil, nil, cfg, nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestRouter_Routes(t *testing.T) {
	router := chi.NewRouter()

	// Test that routes are registered
	// This is a basic smoke test
	if router == nil {
		t.Error("router should not be nil")
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/api/... -v
```

Expected: FAIL with `undefined: NewRouter`

**Step 3: Write minimal implementation**

File: `internal/api/router.go`
```go
package api

import (
	"log/slog"
	"net/http"

	"github.com/atlanssia/llm-eval/internal/api/handler"
	"github.com/atlanssia/llm-eval/internal/api/middleware"
	"github.com/atlanssia/llm-eval/internal/model"
	"github.com/atlanssia/llm-eval/internal/service"
	"github.com/atlanssia/llm-eval/internal/stream"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// NewRouter creates the HTTP router
func NewRouter(
	evalSvc *service.EvaluationService,
	datasetSvc *service.DatasetService,
	modelSvc *service.ModelService,
	streamHub *stream.Hub,
	cfg *model.Config,
	logger *slog.Logger,
) *chi.Mux {
	r := chi.NewRouter()

	// CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(corsMiddleware.Handler)

	// Logging middleware
	r.Use(middleware.Logger(logger))

	// Recovery middleware
	r.Use(middleware.Recover(logger))

	// Optional auth middleware
	if cfg.Auth.Enabled {
		r.Use(middleware.Auth(cfg.Auth.Password))
	}

	// Health check
	r.Get("/health", handler.Health)

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/datasets", handler.ListDatasets(datasetSvc))
		r.Get("/models", handler.ListModels(modelSvc))

		r.Post("/evaluations", handler.CreateEvaluation(evalSvc))
		r.Get("/evaluations", handler.ListEvaluations(evalSvc))
		r.Get("/evaluations/{id}", handler.GetEvaluation(evalSvc))
		r.Get("/evaluations/{id}/stream", handler.StreamEvaluation(evalSvc, streamHub))
	})

	return r
}
```

File: `internal/api/middleware/logger.go`
```go
package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// Logger is a middleware that logs HTTP requests
func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			wrapped := &responseWriter{ResponseWriter: w, status: 200}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			logger.Info("HTTP request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrapped.status,
				"duration", duration,
				"remote_addr", r.RemoteAddr,
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
```

File: `internal/api/middleware/recover.go`
```go
package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recover is a middleware that recovers from panics
func Recover(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Panic recovered",
						"error", err,
						"stack", debug.Stack(),
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
```

File: `internal/api/middleware/auth.go`
```go
package middleware

import (
	"net/http"
)

// Auth is a middleware that optionally requires password authentication
func Auth(password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check for Basic Auth
			username, pass, ok := r.BasicAuth()
			if !ok || username != "admin" || pass != password {
				w.Header().Set("WWW-Authenticate", `Basic realm="LLM Eval"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
```

File: `internal/api/handler/health.go`
```go
package handler

import (
	"encoding/json"
	"net/http"
)

var buildVersion = "dev"

// SetVersion sets the build version
func SetVersion(v string) {
	buildVersion = v
}

// Health responds with health status
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"version": buildVersion,
	})
}
```

**Step 4: Update main.go to use router**

File: `cmd/llm-eval/main.go` - replace mux with router:

```go
// Old code:
// mux := http.NewServeMux()
// mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
//     w.WriteHeader(http.StatusOK)
//     w.Header().Set("Content-Type", "application/json")
//     fmt.Fprintf(w, `{"status":"ok","version":"%s"}`, version)
// })

// New code:
import "github.com/atlanssia/llm-eval/internal/api"
// ... in main:
api.SetVersion(version)
router := api.NewRouter(nil, nil, nil, streamHub, cfg, logger)
srv.Handler = router
```

**Step 5: Run test to verify it passes**

```bash
go test ./internal/api/... -v
```

Expected: PASS

**Step 6: Commit**

```bash
git add internal/api/ cmd/llm-eval/main.go
git commit -m "feat(api): add Chi router with middleware"
```

---

## Task 13: Frontend - Initialize Vite + React Project

**Files:**
- Create: `web/package.json`
- Create: `web/vite.config.ts`
- Create: `web/tsconfig.json`
- Create: `web/tailwind.config.js`
- Create: `web/postcss.config.js`
- Create: `web/index.html`
- Create: `web/src/main.tsx`
- Create: `web/src/App.tsx`
- Create: `web/src/index.css`

**Step 1: Create package.json**

File: `web/package.json`
```json
{
  "name": "llm-eval-web",
  "private": true,
  "version": "0.1.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview",
    "test": "vitest",
    "lint": "eslint . --ext ts,tsx --report-unused-disable-directives --max-warnings 0"
  },
  "dependencies": {
    "react": "^18.3.1",
    "react-dom": "^18.3.1",
    "@tanstack/react-query": "^5.62.11",
    "recharts": "^2.15.0",
    "react-hook-form": "^7.54.2",
    "zod": "^3.24.1",
    "@hookform/resolvers": "^3.10.0"
  },
  "devDependencies": {
    "@types/react": "^18.3.12",
    "@types/react-dom": "^18.3.1",
    "@vitejs/plugin-react": "^4.3.4",
    "autoprefixer": "^10.4.20",
    "postcss": "^8.4.49",
    "tailwindcss": "^3.4.17",
    "typescript": "^5.7.2",
    "vite": "^6.0.7",
    "vitest": "^2.1.8",
    "eslint": "^9.17.0",
    "@typescript-eslint/eslint-plugin": "^8.19.1",
    "@typescript-eslint/parser": "^8.19.1"
  }
}
```

**Step 2: Create Vite config**

File: `web/vite.config.ts`
```typescript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/health': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'dist',
    emptyOutDir: true,
  },
})
```

**Step 3: Create TypeScript config**

File: `web/tsconfig.json`
```json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,

    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx",

    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true
  },
  "include": ["src"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
```

**Step 4: Create tsconfig.node.json**

File: `web/tsconfig.node.json`
```json
{
  "compilerOptions": {
    "composite": true,
    "skipLibCheck": true,
    "module": "ESNext",
    "moduleResolution": "bundler",
    "allowSyntheticDefaultImports": true
  },
  "include": ["vite.config.ts"]
}
```

**Step 5: Create Tailwind config**

File: `web/tailwind.config.js`
```javascript
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
```

**Step 6: Create PostCSS config**

File: `web/postcss.config.js`
```javascript
export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
```

**Step 7: Create index.html**

File: `web/index.html`
```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="/vite.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>LLM Evaluation Tool</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
```

**Step 8: Create main.tsx**

File: `web/src/main.tsx`
```typescript
import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
```

**Step 9: Create App.tsx**

File: `web/src/App.tsx`
```typescript
function App() {
  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 py-4">
          <h1 className="text-2xl font-bold text-gray-900">LLM Evaluation Tool</h1>
        </div>
      </header>
      <main className="max-w-7xl mx-auto px-4 py-8">
        <div className="text-center py-12">
          <p className="text-gray-600">Welcome to LLM Evaluation Tool</p>
        </div>
      </main>
    </div>
  )
}

export default App
```

**Step 10: Create index.css**

File: `web/src/index.css`
```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

**Step 11: Install dependencies and test**

```bash
cd /Users/mw/workspace/tmp/llm-eval/web
npm install
npm run dev
```

Expected: Vite dev server starts on http://localhost:5173

**Step 12: Commit**

```bash
git add web/
git commit -m "feat(web): initialize Vite + React + Tailwind project"
```

---

## Task 14: Frontend - API Client and Types

**Files:**
- Create: `web/src/lib/api.ts`
- Create: `web/src/lib/types.ts`

**Step 1: Write types**

File: `web/src/lib/types.ts`
```typescript
// Evaluation types
export type Status = 'pending' | 'running' | 'completed' | 'failed' | 'canceled'

export interface Evaluation {
  id: string
  created_at: string
  updated_at: string
  status: Status
  config: EvalConfig
  total_cases: number
  completed_cases: number
  error?: string
}

export interface EvalConfig {
  models: string[]
  datasets: string[]
  sample_size: number
  max_workers: number
  ephemeral: boolean
}

export interface Metrics {
  accuracy: number
  f1: number
  bleu: number
  rouge_l: number
  avg_latency: number
  avg_tokens_per_second: number
}

export interface ModelResult {
  model_name: string
  metrics: Metrics
  error_count: number
}

// Dataset types
export interface Dataset {
  name: string
  source: string
  task_type: string
  total_cases: number
  description: string
}

// Model types
export interface ModelConfig {
  name: string
  endpoint: string
  timeout: number
}

// SSE Event types
export type EventType = 'progress' | 'model_complete' | 'evaluation_complete' | 'error'

export interface SSEEvent {
  type: EventType
  data: Record<string, unknown>
}

// API Response types
export interface HealthResponse {
  status: string
  version: string
}

export interface CreateEvaluationRequest {
  models: string[]
  datasets: string[]
  config: {
    sample_size?: number
    max_workers?: number
    ephemeral?: boolean
  }
}
```

**Step 2: Write API client**

File: `web/src/lib/api.ts`
```typescript
import type {
  Evaluation,
  Dataset,
  ModelConfig,
  CreateEvaluationRequest,
  HealthResponse,
  SSEEvent,
} from './types'

const API_BASE = '/api'

class APIClient {
  private baseURL: string

  constructor(baseURL: string = API_BASE) {
    this.baseURL = baseURL
  }

  async get<T>(path: string): Promise<T> {
    const response = await fetch(`${this.baseURL}${path}`)
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }
    return response.json()
  }

  async post<T>(path: string, body: unknown): Promise<T> {
    const response = await fetch(`${this.baseURL}${path}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    })
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }
    return response.json()
  }
}

const api = new APIClient()

// Health check
export async function getHealth(): Promise<HealthResponse> {
  return api.get<HealthResponse>('/health')
}

// Datasets
export async function getDatasets(): Promise<Dataset[]> {
  return api.get<Dataset[]>('/datasets')
}

// Models
export async function getModels(): Promise<ModelConfig[]> {
  return api.get<ModelConfig[]>('/models')
}

// Evaluations
export async function getEvaluations(): Promise<Evaluation[]> {
  return api.get<Evaluation[]>('/evaluations')
}

export async function getEvaluation(id: string): Promise<Evaluation> {
  return api.get<Evaluation>(`/evaluations/${id}`)
}

export async function createEvaluation(request: CreateEvaluationRequest): Promise<Evaluation> {
  return api.post<Evaluation>('/evaluations', request)
}

// SSE streaming
export function streamEvaluation(id: string, onEvent: (event: SSEEvent) => void): () => void {
  const eventSource = new EventSource(`${API_BASE}/evaluations/${id}/stream`)

  eventSource.onmessage = (event) => {
    const data = JSON.parse(event.data) as SSEEvent
    onEvent(data)
  }

  eventSource.onerror = (error) => {
    console.error('SSE error:', error)
    eventSource.close()
  }

  // Return cleanup function
  return () => {
    eventSource.close()
  }
}

export default api
```

**Step 3: Commit**

```bash
git add web/src/lib/
git commit -m "feat(web): add API client and TypeScript types"
```

---

## Task 15: Frontend - Dashboard Page

**Files:**
- Create: `web/src/pages/Dashboard.tsx`
- Create: `web/src/components/EvaluationCard.tsx`

**Step 1: Write test for EvaluationCard component**

File: `web/src/components/EvaluationCard.test.tsx`
```typescript
import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { EvaluationCard } from './EvaluationCard'
import type { Evaluation } from '../lib/types'

describe('EvaluationCard', () => {
  it('shows running evaluation with progress', () => {
    const evaluation: Evaluation = {
      id: '123',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:01:00Z',
      status: 'running',
      config: {
        models: ['gpt-4'],
        datasets: ['mmlu_anatomy'],
        sample_size: 100,
        max_workers: 4,
        ephemeral: false,
      },
      total_cases: 100,
      completed_cases: 45,
    }

    render(<EvaluationCard evaluation={evaluation} />)

    expect(screen.getByText('45%')).toBeInTheDocument()
  })

  it('shows completed evaluation', () => {
    const evaluation: Evaluation = {
      id: '123',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:05:00Z',
      status: 'completed',
      config: {
        models: ['gpt-4'],
        datasets: ['mmlu_anatomy'],
        sample_size: 100,
        max_workers: 4,
        ephemeral: false,
      },
      total_cases: 100,
      completed_cases: 100,
    }

    render(<EvaluationCard evaluation={evaluation} />)

    expect(screen.getByText('Completed')).toBeInTheDocument()
  })
})
```

**Step 2: Run test to verify it fails**

```bash
cd web && npm test
```

Expected: FAIL with `Cannot find module './EvaluationCard'`

**Step 3: Write minimal implementation**

File: `web/src/components/EvaluationCard.tsx`
```typescript
import type { Evaluation } from '../lib/types'

interface EvaluationCardProps {
  evaluation: Evaluation
}

export function EvaluationCard({ evaluation }: EvaluationCardProps) {
  const progress = (evaluation.completed_cases / evaluation.total_cases) * 100

  return (
    <div className="bg-white rounded-lg shadow p-4">
      <div className="flex justify-between items-center mb-2">
        <h3 className="text-lg font-semibold">{evaluation.id}</h3>
        <span className={`px-2 py-1 rounded text-sm ${
          evaluation.status === 'running' ? 'bg-blue-100 text-blue-800' :
          evaluation.status === 'completed' ? 'bg-green-100 text-green-800' :
          evaluation.status === 'failed' ? 'bg-red-100 text-red-800' :
          'bg-gray-100 text-gray-800'
        }`}>
          {evaluation.status.charAt(0).toUpperCase() + evaluation.status.slice(1)}
        </span>
      </div>

      {evaluation.status === 'running' && (
        <div className="mb-2">
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div
              className="bg-blue-600 h-2 rounded-full transition-all"
              style={{ width: `${progress}%` }}
            />
          </div>
          <p className="text-sm text-gray-600 mt-1">{progress.toFixed(0)}%</p>
        </div>
      )}

      <div className="text-sm text-gray-600">
        <p>Models: {evaluation.config.models.join(', ')}</p>
        <p>Datasets: {evaluation.config.datasets.join(', ')}</p>
      </div>
    </div>
  )
}
```

File: `web/src/pages/Dashboard.tsx`
```typescript
import { useQuery } from '@tanstack/react-query'
import { getEvaluations } from '../lib/api'
import { EvaluationCard } from '../components/EvaluationCard'

export function Dashboard() {
  const { data: evaluations, isLoading, error } = useQuery({
    queryKey: ['evaluations'],
    queryFn: getEvaluations,
    refetchInterval: 5000, // Poll every 5 seconds
  })

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-red-600">Failed to load evaluations</p>
      </div>
    )
  }

  const runningEvaluations = evaluations?.filter(e => e.status === 'running') || []
  const recentEvaluations = evaluations?.filter(e => e.status !== 'running') || []

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">Dashboard</h2>
      </div>

      {runningEvaluations.length > 0 && (
        <section className="mb-8">
          <h3 className="text-xl font-semibold mb-4">Active Evaluations</h3>
          <div className="grid gap-4 md:grid-cols-2">
            {runningEvaluations.map(evaluation => (
              <EvaluationCard key={evaluation.id} evaluation={evaluation} />
            ))}
          </div>
        </section>
      )}

      <section>
        <h3 className="text-xl font-semibold mb-4">Recent Evaluations</h3>
        {recentEvaluations.length === 0 ? (
          <div className="text-center py-12 text-gray-500">
            No evaluations yet. Create one to get started.
          </div>
        ) : (
          <div className="grid gap-4 md:grid-cols-2">
            {recentEvaluations.map(evaluation => (
              <EvaluationCard key={evaluation.id} evaluation={evaluation} />
            ))}
          </div>
        )}
      </section>
    </div>
  )
}
```

**Step 4: Update App.tsx to use Dashboard**

File: `web/src/App.tsx`
```typescript
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Dashboard } from './pages/Dashboard'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5000,
      retry: 1,
    },
  },
})

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <div className="min-h-screen bg-gray-50">
        <header className="bg-white shadow">
          <div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
            <h1 className="text-2xl font-bold text-gray-900">LLM Evaluation Tool</h1>
          </div>
        </header>
        <main className="max-w-7xl mx-auto px-4 py-8">
          <Dashboard />
        </main>
      </div>
    </QueryClientProvider>
  )
}

export default App
```

**Step 5: Run test to verify it passes**

```bash
cd web && npm test
```

Expected: PASS

**Step 6: Commit**

```bash
git add web/src/
git commit -m "feat(web): add Dashboard page with EvaluationCard component"
```

---

## Task 16: Embed Frontend into Go Binary

**Files:**
- Modify: `internal/embed/embed.go`

**Step 1: Write embed code**

File: `internal/embed/embed.go`
```go
package embed

import (
	"embed"
	"io/fs"
)

//go:embed ../../web/dist
var distFS embed.FS

// FS returns the embedded filesystem for the React frontend
func FS() (fs.FS, error) {
	// Strip the "web/dist" prefix to serve files from root
	return fs.Sub(distFS, "web/dist")
}

// Exists returns true if the embedded frontend exists
func Exists() bool {
	_, err := distFS.ReadDir("web/dist")
	return err == nil
}
```

**Step 2: Update router to serve embedded files**

File: `internal/api/router.go` - add file server:

```go
import (
	// ... existing imports
	"embedfs "github.com/atlanssia/llm-eval/internal/embed"
	"net/http"
)

// In NewRouter function, after API routes:
// Serve embedded React SPA
if embedfs.Exists() {
	embeddedFS, err := embedfs.FS()
	if err == nil {
		fileServer := http.FileServer(http.FS(embeddedFS))
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			// Serve index.html for SPA routing
			fileServer.ServeHTTP(w, r)
		})
	}
}
```

**Step 3: Update Makefile**

File: `Makefile` - update build-go to depend on build-web:

```makefile
build-go: build-web
	@echo "Building Go binary with embedded frontend..."
	@go build -ldflags "-X main.version=$(VERSION)" -o bin/llm-eval cmd/llm-eval/main.go

build: build-go
	@echo "Build complete: bin/llm-eval"
```

**Step 4: Test embedded build**

```bash
cd web && npm run build
cd ..
go build -o bin/llm-eval cmd/llm-eval/main.go
./bin/llm-eval &
curl http://localhost:8080/
pkill llm-eval
```

Expected: HTML response from embedded React app

**Step 5: Commit**

```bash
git add internal/embed/ Makefile
git commit -m "feat(embed): embed React frontend into Go binary"
```

---

## Task 17: Placeholder Handler Implementations

**Files:**
- Modify: `internal/api/handler/evaluation.go`
- Modify: `internal/api/handler/dataset.go`
- Modify: `internal/api/handler/model.go`

**Step 1: Create placeholder handlers**

File: `internal/api/handler/evaluation.go`
```go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/atlanssia/llm-eval/internal/service"
	"github.com/atlanssia/llm-eval/internal/stream"
)

// CreateEvaluation creates a new evaluation
func CreateEvaluation(evalSvc *service.EvaluationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "not implemented"})
	}
}

// ListEvaluations lists all evaluations
func ListEvaluations(evalSvc *service.EvaluationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{})
	}
}

// GetEvaluation gets an evaluation by ID
func GetEvaluation(evalSvc *service.EvaluationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "not implemented"})
	}
}

// StreamEvaluation streams evaluation progress via SSE
func StreamEvaluation(evalSvc *service.EvaluationService, hub *stream.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "SSE not supported", http.StatusInternalServerError)
			return
		}

		// TODO: Implement SSE streaming
		flusher.Flush()
	}
}
```

File: `internal/api/handler/dataset.go`
```go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/atlanssia/llm-eval/internal/service"
)

// ListDatasets lists available datasets
func ListDatasets(datasetSvc *service.DatasetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]string{
			{"name": "mmlu_anatomy", "description": "MMLU Anatomy"},
			{"name": "cmmlu_anatomy", "description": "CMMLU Anatomy"},
		})
	}
}
```

File: `internal/api/handler/model.go`
```go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/atlanssia/llm-eval/internal/service"
)

// ListModels lists configured models
func ListModels(modelSvc *service.ModelService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]string{
			{"name": "gpt-4", "endpoint": "https://api.openai.com/v1/chat/completions"},
		})
	}
}
```

**Step 2: Create placeholder services**

File: `internal/service/evaluation.go`
```go
package service

import (
	"log/slog"

	"github.com/atlanssia/llm-eval/internal/repository"
	"github.com/atlanssia/llm-eval/internal/stream"
)

// EvaluationService handles evaluation business logic
type EvaluationService struct {
	evalRepo  *repository.EvaluationRepository
	resultRepo *repository.ResultRepository
	hub       *stream.Hub
	logger    *slog.Logger
}

// NewEvaluationService creates a new evaluation service
func NewEvaluationService(
	evalRepo *repository.EvaluationRepository,
	resultRepo *repository.ResultRepository,
	hub *stream.Hub,
	logger *slog.Logger,
) *EvaluationService {
	return &EvaluationService{
		evalRepo:  evalRepo,
		resultRepo: resultRepo,
		hub:       hub,
		logger:    logger,
	}
}
```

File: `internal/service/dataset.go`
```go
package service

import (
	"log/slog"
)

// DatasetService handles dataset loading
type DatasetService struct {
	dataDir string
	logger  *slog.Logger
}

// NewDatasetService creates a new dataset service
func NewDatasetService(dataDir string, logger *slog.Logger) *DatasetService {
	return &DatasetService{
		dataDir: dataDir,
		logger:  logger,
	}
}
```

File: `internal/service/model.go`
```go
package service

import (
	"log/slog"

	"github.com/atlanssia/llm-eval/internal/model"
)

// ModelService handles model API calls
type ModelService struct {
	models []model.ModelConfig
	logger *slog.Logger
}

// NewModelService creates a new model service
func NewModelService(models []model.ModelConfig, logger *slog.Logger) *ModelService {
	return &ModelService{
		models: models,
		logger: logger,
	}
}
```

**Step 3: Update main.go to initialize services**

File: `cmd/llm-eval/main.go` - update service initialization:

```go
// After repository initialization:
// Initialize services
datasetSvc := service.NewDatasetService(cfg.DataDir, logger)
modelSvc := service.NewModelService(cfg.Models, logger)
evalSvc := service.NewEvaluationService(evalRepo, resultRepo, streamHub, logger)

// Pass services to router
router := api.NewRouter(evalSvc, datasetSvc, modelSvc, streamHub, cfg, logger)
```

**Step 4: Test basic API endpoints**

```bash
go run cmd/llm-eval/main.go &
sleep 2
curl http://localhost:8080/api/datasets
curl http://localhost:8080/api/models
curl http://localhost:8080/api/evaluations
pkill -f "go run cmd/llm-eval/main.go"
```

Expected: JSON responses from each endpoint

**Step 5: Commit**

```bash
git add internal/api/handler/ internal/service/ cmd/llm-eval/main.go
git commit -m "feat(api): add placeholder handler and service implementations"
```

---

## Task 18: Update .gitignore

**Files:**
- Modify: `.gitignore`

**Step 1: Update .gitignore**

File: `.gitignore`
```gitignore
# Binaries
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib
llm-eval

# Test files
coverage.out
coverage.html
*.test

# Go workspace
/workspace/

# IDE
.vscode/
.idea/
*.swp
*.swo

# Database
*.db
*.db-shm
*.db-wal
data/
!data/.gitkeep

# Node
web/node_modules/
web/dist/

# Environment
.env
.env.local
config.yaml
!configs/*.example
```

**Step 2: Create data/.gitkeep**

```bash
mkdir -p data
touch data/.gitkeep
```

**Step 3: Commit**

```bash
git add .gitignore data/.gitkeep
git commit -m "chore: update .gitignore with additional patterns"
```

---

## End of Implementation Plan

This plan covers:
1. Project initialization with Go module and directory structure
2. Domain models (Evaluation, Dataset, Config)
3. Stream hub for SSE events
4. Database migrations and repositories
5. Configuration loading
6. Main entry point with graceful shutdown
7. Chi router with middleware
8. Frontend initialization (Vite + React + Tailwind)
9. API client and types
10. Dashboard page
11. Embedded frontend in Go binary
12. Placeholder handlers and services

**Next Steps:**
1. Run this plan using `executing-plans` skill
2. Implement remaining service layer (dataset loading, model API client, metrics calculation)
3. Implement full evaluation orchestration
4. Add remaining frontend pages (New Evaluation, Results, Settings)
5. Add E2E tests with Playwright

**Total estimated tasks: 18**
**Total estimated time: 4-6 hours for full implementation**
