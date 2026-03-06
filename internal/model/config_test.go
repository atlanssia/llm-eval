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
