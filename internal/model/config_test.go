package model

import (
	"os"
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
				Timeout:  30,
			},
			wantErr: false,
		},
		{
			name: "valid config with timeout and max retries",
			config: ModelConfig{
				Name:       "gpt-4",
				Endpoint:   "https://api.openai.com/v1/chat/completions",
				APIKey:     "sk-test",
				Timeout:    30,
				MaxRetries: 3,
			},
			wantErr: false,
		},
		{
			name: "valid config with zero max retries",
			config: ModelConfig{
				Name:       "gpt-4",
				Endpoint:   "https://api.openai.com/v1/chat/completions",
				APIKey:     "sk-test",
				Timeout:    30,
				MaxRetries: 0,
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
		{
			name: "invalid timeout - zero",
			config: ModelConfig{
				Name:     "gpt-4",
				Endpoint: "https://api.openai.com/v1/chat/completions",
				APIKey:   "sk-test",
				Timeout:  0,
			},
			wantErr: true,
		},
		{
			name: "invalid timeout - negative",
			config: ModelConfig{
				Name:     "gpt-4",
				Endpoint: "https://api.openai.com/v1/chat/completions",
				APIKey:   "sk-test",
				Timeout:  -1,
			},
			wantErr: true,
		},
		{
			name: "invalid max retries - negative",
			config: ModelConfig{
				Name:       "gpt-4",
				Endpoint:   "https://api.openai.com/v1/chat/completions",
				APIKey:     "sk-test",
				Timeout:    30,
				MaxRetries: -1,
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

func TestModelConfig_GetAPIKey(t *testing.T) {
	// Setup test environment variable
	os.Setenv("TEST_API_KEY", "sk-secret-123")
	defer os.Unsetenv("TEST_API_KEY")

	tests := []struct {
		name     string
		config   ModelConfig
		expected string
	}{
		{
			name: "plain API key",
			config: ModelConfig{
				APIKey: "sk-test-key",
			},
			expected: "sk-test-key",
		},
		{
			name: "environment variable expansion",
			config: ModelConfig{
				APIKey: "${TEST_API_KEY}",
			},
			expected: "sk-secret-123",
		},
		{
			name: "empty API key",
			config: ModelConfig{
				APIKey: "",
			},
			expected: "",
		},
		{
			name: "non-existent environment variable",
			config: ModelConfig{
				APIKey: "${NON_EXISTENT_VAR}",
			},
			expected: "",
		},
		{
			name: "malformed environment variable - missing closing brace",
			config: ModelConfig{
				APIKey: "${TEST_API_KEY",
			},
			expected: "${TEST_API_KEY",
		},
		{
			name: "malformed environment variable - missing opening brace",
			config: ModelConfig{
				APIKey: "TEST_API_KEY}",
			},
			expected: "TEST_API_KEY}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetAPIKey()
			if result != tt.expected {
				t.Errorf("GetAPIKey() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config with empty models slice",
			config: Config{
				Server: ServerConfig{
					Addr: "0.0.0.0:8080",
				},
				Database: DatabaseConfig{
					Path: "./data/test.db",
				},
				Models: []ModelConfig{},
			},
			wantErr: false,
		},
		{
			name: "valid config with multiple valid models",
			config: Config{
				Server: ServerConfig{
					Addr: "0.0.0.0:8080",
				},
				Database: DatabaseConfig{
					Path: "./data/test.db",
				},
				Models: []ModelConfig{
					{
						Name:     "gpt-4",
						Endpoint: "https://api.openai.com/v1/chat/completions",
						APIKey:   "sk-test-1",
						Timeout:  30,
					},
					{
						Name:     "claude-3",
						Endpoint: "https://api.anthropic.com/v1/messages",
						APIKey:   "sk-test-2",
						Timeout:  60,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid config with one invalid model",
			config: Config{
				Server: ServerConfig{
					Addr: "0.0.0.0:8080",
				},
				Database: DatabaseConfig{
					Path: "./data/test.db",
				},
				Models: []ModelConfig{
					{
						Name:     "gpt-4",
						Endpoint: "https://api.openai.com/v1/chat/completions",
						APIKey:   "sk-test-1",
						Timeout:  30,
					},
					{
						Name:     "", // Invalid: missing name
						Endpoint: "https://api.anthropic.com/v1/messages",
						APIKey:   "sk-test-2",
						Timeout:  60,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid config with model having negative timeout",
			config: Config{
				Server: ServerConfig{
					Addr: "0.0.0.0:8080",
				},
				Database: DatabaseConfig{
					Path: "./data/test.db",
				},
				Models: []ModelConfig{
					{
						Name:     "gpt-4",
						Endpoint: "https://api.openai.com/v1/chat/completions",
						APIKey:   "sk-test-1",
						Timeout:  -5, // Invalid: negative timeout
					},
				},
			},
			wantErr: true,
		},
		{
			name: "nil models slice",
			config: Config{
				Server: ServerConfig{
					Addr: "0.0.0.0:8080",
				},
				Database: DatabaseConfig{
					Path: "./data/test.db",
				},
				Models: nil,
			},
			wantErr: false,
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
