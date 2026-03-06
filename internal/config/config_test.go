package config

import (
	"os"
	"path/filepath"
	"testing"
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
