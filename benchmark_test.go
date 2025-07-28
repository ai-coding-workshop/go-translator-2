package main

import (
	"context"
	"testing"
	"time"
	"translator-service/internal/config"
	"translator-service/internal/models"
	"translator-service/internal/services"
)

func BenchmarkTranslatorService(b *testing.B) {
	// Create a test config
	cfg := &config.Config{
		ServerPort: "8080",
		Timeout:    30,
	}

	// Create translator service
	ts := services.NewTranslatorService(cfg)

	// Create a test request
	req := &models.TranslationRequest{
		Text:  "Hello, world! This is a benchmark test.",
		Model: "gpt-3.5",
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Perform translation
		_, err := ts.Translate(ctx, req)
		if err != nil {
			b.Fatalf("Translation failed: %v", err)
		}
	}
}
