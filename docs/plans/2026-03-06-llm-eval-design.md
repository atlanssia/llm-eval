# LLM Evaluation Web Tool - Design Document

**Date:** 2026-03-06
**Author:** Claude Code
**Status:** Approved

---

## 1. Overview

A web-based LLM evaluation tool built with Go + React + Tailwind + Vite. The frontend React application is embedded into a single Go binary for easy distribution. The tool provides real-time evaluation monitoring with rich dashboard visualization.

**Key Requirements:**
- Interactive real-time monitoring + rich dashboard visualization
- Both config file upload AND visual form builder
- Hybrid persistence (disk storage with ephemeral option)
- LAN access with optional password protection

---

## 2. Project Structure

```
llm-eval/
├── cmd/
│   └── llm-eval/
│       └── main.go                 # Entry point with graceful shutdown
├── internal/
│   ├── api/
│   │   ├── handler/                # HTTP handlers (thin, delegate to services)
│   │   │   ├── evaluation.go       # POST /api/evaluations, GET /api/evaluations/:id/stream
│   │   │   ├── dataset.go          # GET /api/datasets
│   │   │   ├── model.go            # GET /api/models
│   │   │   └── health.go           # GET /health
│   │   ├── middleware/             # Chi middleware
│   │   │   ├── auth.go             # Optional password via BasicAuth
│   │   │   ├── cors.go             # CORS handling
│   │   │   ├── logger.go           # Request logging
│   │   │   └── recover.go          # Panic recovery
│   │   ├── dto/                    # Request/Response DTOs
│   │   │   ├── evaluation.go
│   │   │   ├── dataset.go
│   │   │   └── common.go           # Pagination, error responses
│   │   └── router.go               # Chi router setup, route registration
│   ├── service/
│   │   ├── evaluation.go           # Core evaluation orchestration
│   │   ├── dataset.go              # Dataset loading (MMLU/CMMLU/CSV)
│   │   ├── model.go                # LLM API client with retry logic
│   │   ├── metrics.go              # Accuracy, F1, BLEU, ROUGE calculation
│   │   └── stream.go               # SSE event broadcasting
│   ├── repository/
│   │   ├── evaluation.go           # SQLite CRUD for evaluations
│   │   ├── result.go               # SQLite CRUD for results
│   │   └── migrations.go           # Database schema setup
│   ├── model/
│   │   ├── evaluation.go           # Domain entity: Evaluation
│   │   ├── dataset.go              # Domain entity: Dataset, Case
│   │   ├── result.go               # Domain entity: Result, ModelResult
│   │   └── config.go               # Domain entity: ModelConfig, EvalConfig
│   └── stream/
│   │   ├── events.go               # SSE event types
│   │   └── hub.go                  # Event broadcasting hub
│   └── embed/
│       └── embed.go                # //go:embed for React dist
├── web/                            # React frontend (Vite project)
│   ├── src/
│   │   ├── components/             # React components
│   │   ├── pages/                  # Route pages
│   │   ├── lib/                    # API client, types
│   │   └── main.tsx
│   ├── index.html
│   ├── vite.config.ts
│   └── tailwind.config.js
├── configs/
│   └── models.yaml.example         # Example model config
├── migrations/
│   └── 001_init.sql                # SQLite schema
├── Makefile                        # Build automation
├── go.mod
└── go.sum
```

---

## 3. Technology Stack

| Layer | Technology |
|-------|-----------|
| **Go Version** | 1.26.1 |
| **HTTP Framework** | `github.com/go-chi/chi/v5` |
| **Database** | SQLite (`modernc.org/sqlite` - CGo-free) |
| **Config** | YAML (`gopkg.in/yaml.v3`) |
| **Frontend Framework** | React 18 + TypeScript |
| **Build Tool** | Vite |
| **Styling** | Tailwind CSS |
| **Components** | shadcn/ui + Radix UI |
| **Charts** | Recharts |
| **Forms** | React Hook Form + Zod |
| **State** | TanStack Query (React Query) |
| **Routing** | TanStack Router |
| **Real-time** | Server-Sent Events (SSE) |

---

## 4. Architecture

### HTTP API Routes

```
GET  /api/datasets                → List available datasets
GET  /api/models                  → List configured models
POST /api/evaluations             → Start new evaluation
GET  /api/evaluations             → List all evaluations
GET  /api/evaluations/:id         → Get evaluation details
GET  /api/evaluations/:id/stream  → SSE: real-time progress
GET  /api/results/:id/export      → Export results (JSON/CSV)
GET  /health                      → Health check
GET  /*                           → React SPA (fallback)
```

### Request Flow: Start Evaluation

1. User fills form / uploads config (React)
2. POST /api/evaluations { models: [...], datasets: [...], config: {...} }
3. EvaluationService.Create()
   - Validate config
   - Load datasets from disk
   - Create evaluation record (SQLite)
   - Return evaluation ID immediately
4. Background goroutine starts:
   - For each model (parallel):
     - For each dataset case (parallel):
       - Call LLM API
       - Record prediction + latency
       - Calculate metrics
       - Emit SSE event: { type: "progress", model, dataset, current, total }
     - Emit SSE event: { type: "model_complete", model, metrics }
   - Emit SSE event: { type: "evaluation_complete", summary }
5. Frontend receives SSE events → Update UI in real-time

### Graceful Shutdown

- Catch SIGTERM/SIGINT signals
- 30-second shutdown timeout
- Stop accepting new connections
- Close stream hub (disconnects SSE clients)
- Close database connections
- Clean resource cleanup via defer

---

## 5. Frontend Design

### Page Structure

- **/** (Dashboard)
  - Active evaluations with live progress
  - Recent history table

- **/evaluations/new**
  - Step-by-step wizard: models → datasets → config → review

- **/evaluations/:id**
  - Real-time evaluation details
  - Streaming log (SSE)
  - Per-model progress tabs

- **/results/:id**
  - Metrics summary cards
  - Comparison charts
  - Sortable/filterable results table

- **/settings**
  - Model configuration
  - Data source paths
  - Security settings

### Component Hierarchy

```
AppLayout
├── Header
│   ├── Logo
│   ├── NavLinks
│   └── UserMenu
└── MainContent
    ├── DashboardPage
    │   ├── ActiveEvaluationCard[]
    │   └── RecentHistoryTable
    ├── NewEvaluationPage
    │   ├── ModelSelector
    │   ├── DatasetSelector
    │   ├── ConfigForm
    │   └── ReviewSummary
    ├── EvaluationDetailPage
    │   ├── EvaluationHeader
    │   ├── StreamingLog
    │   └── ModelProgressTabs
    └── ResultsPage
        ├── MetricsSummaryCards
        ├── ComparisonBarChart
        └── ResultsTable
```

---

## 6. Go Best Practices Applied

1. **No `pkg/` directory** - Use `internal/` for all private code
2. **`main.go` in `cmd/`** - Entry point should be minimal
3. **Interface-based design** - Services depend on interfaces, not concretions
4. **Context propagation** - Pass `context.Context` for cancellation/timeout
5. **Structured logging** - Use `slog` (Go 1.21+)
6. **Error wrapping** - Use `fmt.Errorf` with `%w` for error chains
7. **Graceful shutdown** - Proper signal handling with timeout
8. **Connection pooling** - SQLite pool limits configured
9. **Server timeouts** - Read/Write/Idle timeouts configured
10. **Resource cleanup** - defer Close() on all resources

---

## 7. Testing Strategy (TDD)

### Testing Pyramid

```
┌─────────────────┐
│   E2E Tests     │  ← Playwright (5%)
├─────────────────┤
│  Integration    │  ← API tests with SQLite :memory: (15%)
├─────────────────┤
│   Unit Tests    │  ← Go tests + React Vitest (80%)
└─────────────────┘
```

### TDD Workflow

1. **Red**: Write failing test → `go test ./...` → ❌ FAIL
2. **Green**: Write minimal code → `go test ./...` → ✅ PASS
3. **Refactor**: Clean up → `go test ./...` → ✅ Still PASS
4. **Repeat**: Next feature

### Test Structure

```
internal/
├── api/handler/
│   └── *_test.go              # HTTP handler tests (httptest)
├── service/
│   └── *_test.go              # Service layer tests (mocked repos)
├── repository/
│   └── *_test.go              # Repository tests (SQLite :memory:)
└── model/
    └── *_test.go              # Domain model tests (pure functions)
```

---

## 8. Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_ADDR` | Server bind address | `0.0.0.0:8080` |
| `DATABASE_PATH` | SQLite database path | `./data/llm-eval.db` |
| `DATA_DIR` | Local datasets directory | `./data` |
| `AUTH_ENABLED` | Enable password auth | `false` |
| `AUTH_PASSWORD` | Password for auth | - |

### Config File (YAML)

```yaml
server:
  addr: "0.0.0.0:8080"
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "120s"

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
```

---

## 9. Build & Deployment

### Makefile Targets

```makefile
make dev         # Start Go API + Vite dev server
make build-web   # Build React frontend
make build-go    # Build Go binary with embedded frontend
make build       # Full production build
make test        # Run all tests
make test-go     # Run Go tests with coverage
make test-web    # Run frontend tests
make run         # Run the binary
```

### Single Binary Distribution

The final binary includes:
- Compiled Go code
- Embedded React static assets
- SQLite database (created at runtime)

---

## 10. Package Naming

All packages use the prefix: `github.com/atlanssia/llm-eval`

---

**Next Steps:** Invoke `writing-plans` skill to create detailed implementation plan.
