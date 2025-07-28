package config

import (
	"os"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "Valid config",
			config: &Config{
				ServerPort: "8080",
				Timeout:    30,
			},
			expectError: false,
		},
		{
			name: "Empty port",
			config: &Config{
				ServerPort: "",
				Timeout:    30,
			},
			expectError: true,
		},
		{
			name: "Invalid port",
			config: &Config{
				ServerPort: "invalid",
				Timeout:    30,
			},
			expectError: true,
		},
		{
			name: "Valid port number",
			config: &Config{
				ServerPort: "9000",
				Timeout:    30,
			},
			expectError: false,
		},
		{
			name: "Zero timeout",
			config: &Config{
				ServerPort: "8080",
				Timeout:    0,
			},
			expectError: true,
		},
		{
			name: "Negative timeout",
			config: &Config{
				ServerPort: "8080",
				Timeout:    -1,
			},
			expectError: true,
		},
		{
			name: "Timeout too large",
			config: &Config{
				ServerPort: "8080",
				Timeout:    301,
			},
			expectError: true,
		},
		{
			name: "Valid large timeout",
			config: &Config{
				ServerPort: "8080",
				Timeout:    300,
			},
			expectError: false,
		},
		{
			name: "Invalid OpenAI endpoint",
			config: &Config{
				ServerPort:     "8080",
				Timeout:        30,
				OpenAIEndpoint: "invalid-url",
			},
			expectError: true,
		},
		{
			name: "Valid OpenAI endpoint",
			config: &Config{
				ServerPort:     "8080",
				Timeout:        30,
				OpenAIEndpoint: "https://api.openai.com/v1",
			},
			expectError: false,
		},
		{
			name: "Invalid Anthropic endpoint",
			config: &Config{
				ServerPort:        "8080",
				Timeout:           30,
				AnthropicEndpoint: "invalid-url",
			},
			expectError: true,
		},
		{
			name: "Valid Anthropic endpoint",
			config: &Config{
				ServerPort:        "8080",
				Timeout:           30,
				AnthropicEndpoint: "https://api.anthropic.com/v1",
			},
			expectError: false,
		},
		{
			name: "OpenAI key without endpoint",
			config: &Config{
				ServerPort: "8080",
				Timeout:    30,
				OpenAIKey:  "test-key",
			},
			expectError: true,
		},
		{
			name: "Anthropic key without endpoint",
			config: &Config{
				ServerPort:   "8080",
				Timeout:      30,
				AnthropicKey: "test-key",
			},
			expectError: true,
		},
		{
			name: "OpenAI key with endpoint",
			config: &Config{
				ServerPort:     "8080",
				Timeout:        30,
				OpenAIEndpoint: "https://api.openai.com/v1",
				OpenAIKey:      "test-key",
			},
			expectError: false,
		},
		{
			name: "Anthropic key with endpoint",
			config: &Config{
				ServerPort:        "8080",
				Timeout:           30,
				AnthropicEndpoint: "https://api.anthropic.com/v1",
				AnthropicKey:      "test-key",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestConfig_HasKeys(t *testing.T) {
	config := &Config{
		OpenAIKey:    "test-openai-key",
		AnthropicKey: "",
	}

	if !config.HasOpenAIKey() {
		t.Errorf("Expected HasOpenAIKey to return true")
	}

	if config.HasAnthropicKey() {
		t.Errorf("Expected HasAnthropicKey to return false")
	}
}

func TestConfig_GetKeysAndEndpoints(t *testing.T) {
	config := &Config{
		OpenAIEndpoint:    "https://custom.openai.com/v1",
		OpenAIKey:         "test-openai-key",
		AnthropicEndpoint: "https://custom.anthropic.com/v1",
		AnthropicKey:      "test-anthropic-key",
	}

	if config.GetOpenAIKey() != "test-openai-key" {
		t.Errorf("Expected GetOpenAIKey to return test-openai-key")
	}

	if config.GetAnthropicKey() != "test-anthropic-key" {
		t.Errorf("Expected GetAnthropicKey to return test-anthropic-key")
	}

	if config.GetOpenAIEndpoint() != "https://custom.openai.com/v1" {
		t.Errorf("Expected GetOpenAIEndpoint to return custom endpoint")
	}

	if config.GetAnthropicEndpoint() != "https://custom.anthropic.com/v1" {
		t.Errorf("Expected GetAnthropicEndpoint to return custom endpoint")
	}

	// Test default endpoints
	defaultConfig := &Config{}
	if defaultConfig.GetOpenAIEndpoint() != "https://api.openai.com/v1" {
		t.Errorf("Expected default OpenAI endpoint")
	}

	if defaultConfig.GetAnthropicEndpoint() != "https://api.anthropic.com/v1" {
		t.Errorf("Expected default Anthropic endpoint")
	}
}

func TestConfig_LoadFromEnv(t *testing.T) {
	// Save original environment variables
	originalPort := os.Getenv("PORT")
	originalOpenAIKey := os.Getenv("OPENAI_API_KEY")
	originalAnthropicKey := os.Getenv("ANTHROPIC_API_KEY")
	originalDebug := os.Getenv("DEBUG")
	originalTimeout := os.Getenv("TIMEOUT")

	// Clean up
	defer func() {
		os.Setenv("PORT", originalPort)
		os.Setenv("OPENAI_API_KEY", originalOpenAIKey)
		os.Setenv("ANTHROPIC_API_KEY", originalAnthropicKey)
		os.Setenv("DEBUG", originalDebug)
		os.Setenv("TIMEOUT", originalTimeout)
	}()

	// Set test environment variables
	os.Setenv("PORT", "9000")
	os.Setenv("OPENAI_API_KEY", "test-env-openai-key")
	os.Setenv("ANTHROPIC_API_KEY", "test-env-anthropic-key")
	os.Setenv("DEBUG", "true")
	os.Setenv("TIMEOUT", "60")

	config := &Config{
		ServerPort:        "8080",
		OpenAIKey:         "",
		AnthropicKey:      "",
		Debug:             false,
		Timeout:           30,
		OpenAIEndpoint:    "https://api.openai.com/v1",
		AnthropicEndpoint: "https://api.anthropic.com/v1",
	}

	config.loadFromEnv()

	if config.ServerPort != "9000" {
		t.Errorf("Expected ServerPort to be 9000, got %s", config.ServerPort)
	}

	if config.OpenAIKey != "test-env-openai-key" {
		t.Errorf("Expected OpenAIKey to be test-env-openai-key, got %s", config.OpenAIKey)
	}

	if config.AnthropicKey != "test-env-anthropic-key" {
		t.Errorf("Expected AnthropicKey to be test-env-anthropic-key, got %s", config.AnthropicKey)
	}

	if !config.Debug {
		t.Errorf("Expected Debug to be true")
	}

	if config.Timeout != 60 {
		t.Errorf("Expected Timeout to be 60, got %d", config.Timeout)
	}
}
