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
	go run cmd/llm-eval/main.go

dev-web:
	@echo "Starting Vite dev server on :5173..."
	cd web && npm run dev

build-web:
	@echo "Building React frontend..."
	cd web && npm run build

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
