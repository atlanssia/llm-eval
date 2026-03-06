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

	// Apply defaults to models
	for i := range cfg.Models {
		if cfg.Models[i].Timeout == 0 {
			cfg.Models[i].Timeout = 60
		}
		if cfg.Models[i].MaxRetries == 0 {
			cfg.Models[i].MaxRetries = 2
		}
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
