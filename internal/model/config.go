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
	if c.Timeout <= 0 {
		return errors.New("model timeout must be greater than 0")
	}
	if c.MaxRetries < 0 {
		return errors.New("model max retries must be greater than or equal to 0")
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
