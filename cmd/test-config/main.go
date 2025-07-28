package main

import (
	"fmt"
	"log"
	"os"
	"translator-service/internal/config"
)

func main() {
	// Set some environment variables for testing
	os.Setenv("OPENAI_API_KEY", "test-openai-key")
	os.Setenv("ANTHROPIC_API_KEY", "test-anthropic-key")
	os.Setenv("PORT", "9090")

	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Print configuration
	fmt.Printf("Server Port: %s\n", cfg.ServerPort)
	fmt.Printf("OpenAI Endpoint: %s\n", cfg.GetOpenAIEndpoint())
	fmt.Printf("Has OpenAI Key: %t\n", cfg.HasOpenAIKey())
	fmt.Printf("Anthropic Endpoint: %s\n", cfg.GetAnthropicEndpoint())
	fmt.Printf("Has Anthropic Key: %t\n", cfg.HasAnthropicKey())
	fmt.Printf("Timeout: %d\n", cfg.Timeout)
	fmt.Printf("Debug: %t\n", cfg.Debug)
}
