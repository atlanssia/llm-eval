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
