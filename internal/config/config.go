package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds application configuration
type Config struct {
	ServerPort        string `yaml:"port"`
	OpenAIEndpoint    string `yaml:"openai_endpoint"`
	OpenAIKey         string `yaml:"openai_key"`
	AnthropicEndpoint string `yaml:"anthropic_endpoint"`
	AnthropicKey      string `yaml:"anthropic_key"`
	Debug             bool   `yaml:"debug"`
	Timeout           int    `yaml:"timeout"`
}

// NewConfig creates a new configuration from environment variables and config file
func NewConfig() (*Config, error) {
	// Parse command line flags
	configFile := flag.String("config", "", "Path to config file")
	flag.Parse()

	// Create default config
	config := &Config{
		ServerPort:        "8080",
		OpenAIEndpoint:    "https://api.openai.com/v1",
		OpenAIKey:         "",
		AnthropicEndpoint: "https://api.anthropic.com/v1",
		AnthropicKey:      "",
		Debug:             false,
		Timeout:           30,
	}

	// Load from config file if specified
	if *configFile != "" {
		if err := config.loadFromFile(*configFile); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}

	// Override with environment variables
	config.loadFromEnv()

	// Validate configuration
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// loadFromFile loads configuration from a YAML file
func (c *Config) loadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var fileConfig struct {
		Server struct {
			Port string `yaml:"port"`
		} `yaml:"server"`
		LLM struct {
			OpenAIEndpoint    string `yaml:"openai_endpoint"`
			OpenAIKey         string `yaml:"openai_key"`
			AnthropicEndpoint string `yaml:"anthropic_endpoint"`
			AnthropicKey      string `yaml:"anthropic_key"`
			Timeout           int    `yaml:"timeout"`
		} `yaml:"llm"`
		Debug bool `yaml:"debug"`
	}

	if err := yaml.Unmarshal(data, &fileConfig); err != nil {
		return err
	}

	// Apply file config
	if fileConfig.Server.Port != "" {
		c.ServerPort = fileConfig.Server.Port
	}
	if fileConfig.LLM.OpenAIEndpoint != "" {
		c.OpenAIEndpoint = fileConfig.LLM.OpenAIEndpoint
	}
	if fileConfig.LLM.OpenAIKey != "" {
		c.OpenAIKey = fileConfig.LLM.OpenAIKey
	}
	if fileConfig.LLM.AnthropicEndpoint != "" {
		c.AnthropicEndpoint = fileConfig.LLM.AnthropicEndpoint
	}
	if fileConfig.LLM.AnthropicKey != "" {
		c.AnthropicKey = fileConfig.LLM.AnthropicKey
	}
	if fileConfig.LLM.Timeout > 0 {
		c.Timeout = fileConfig.LLM.Timeout
	}
	c.Debug = fileConfig.Debug

	return nil
}

// loadFromEnv loads configuration from environment variables
func (c *Config) loadFromEnv() {
	if value := os.Getenv("PORT"); value != "" {
		c.ServerPort = value
	}
	if value := os.Getenv("OPENAI_ENDPOINT"); value != "" {
		c.OpenAIEndpoint = value
	}
	if value := os.Getenv("OPENAI_API_KEY"); value != "" {
		c.OpenAIKey = value
	}
	if value := os.Getenv("ANTHROPIC_ENDPOINT"); value != "" {
		c.AnthropicEndpoint = value
	}
	if value := os.Getenv("ANTHROPIC_API_KEY"); value != "" {
		c.AnthropicKey = value
	}
	if value := os.Getenv("DEBUG"); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			c.Debug = boolValue
		}
	}
	if value := os.Getenv("TIMEOUT"); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			c.Timeout = intValue
		}
	}
}

// validate checks that the configuration is valid
func (c *Config) validate() error {
	// Validate server port
	if c.ServerPort == "" {
		return fmt.Errorf("server port cannot be empty")
	}

	// Validate that server port is a valid number
	if _, err := strconv.Atoi(c.ServerPort); err != nil {
		return fmt.Errorf("server port must be a valid number: %w", err)
	}

	// Validate timeout
	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	// Validate timeout range (reasonable limits)
	if c.Timeout > 300 {
		return fmt.Errorf("timeout is too large (maximum 300 seconds)")
	}

	// Validate endpoints
	if c.OpenAIEndpoint != "" && !strings.HasPrefix(c.OpenAIEndpoint, "http") {
		return fmt.Errorf("openai_endpoint must be a valid URL")
	}
	if c.AnthropicEndpoint != "" && !strings.HasPrefix(c.AnthropicEndpoint, "http") {
		return fmt.Errorf("anthropic_endpoint must be a valid URL")
	}

	// Validate that if API keys are provided, endpoints are also provided
	if c.OpenAIKey != "" && c.OpenAIEndpoint == "" {
		return fmt.Errorf("openai_endpoint must be provided when openai_key is set")
	}
	if c.AnthropicKey != "" && c.AnthropicEndpoint == "" {
		return fmt.Errorf("anthropic_endpoint must be provided when anthropic_key is set")
	}

	return nil
}

// HasOpenAIKey returns true if an OpenAI API key is configured
func (c *Config) HasOpenAIKey() bool {
	return c.OpenAIKey != ""
}

// HasAnthropicKey returns true if an Anthropic API key is configured
func (c *Config) HasAnthropicKey() bool {
	return c.AnthropicKey != ""
}

// GetOpenAIKey returns the OpenAI API key, or an empty string if not configured
func (c *Config) GetOpenAIKey() string {
	return c.OpenAIKey
}

// GetAnthropicKey returns the Anthropic API key, or an empty string if not configured
func (c *Config) GetAnthropicKey() string {
	return c.AnthropicKey
}

// GetOpenAIEndpoint returns the OpenAI endpoint, or the default if not configured
func (c *Config) GetOpenAIEndpoint() string {
	if c.OpenAIEndpoint == "" {
		return "https://api.openai.com/v1"
	}
	return c.OpenAIEndpoint
}

// GetAnthropicEndpoint returns the Anthropic endpoint, or the default if not configured
func (c *Config) GetAnthropicEndpoint() string {
	if c.AnthropicEndpoint == "" {
		return "https://api.anthropic.com/v1"
	}
	return c.AnthropicEndpoint
}
