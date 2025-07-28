package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"translator-service/internal/models"
)

// MockTranslator is a mock implementation of the Translator interface
type MockTranslator struct {
	name string
}

// NewMockTranslator creates a new mock translator
func NewMockTranslator(name string) *MockTranslator {
	return &MockTranslator{name: name}
}

// Translate performs a mock translation
func (mt *MockTranslator) Translate(ctx context.Context, req *models.TranslationRequest) (*models.TranslationResponse, error) {
	// Simulate API call delay
	time.Sleep(500 * time.Millisecond)

	// Check for context cancellation
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Create mock translation
	translation := mt.mockTranslate(req.Text)

	return &models.TranslationResponse{
		Original:    req.Text,
		Translation: translation,
		Model:       mt.name,
	}, nil
}

// Name returns the name of the translator
func (mt *MockTranslator) Name() string {
	return mt.name
}

// SupportsModel returns true if the translator supports the given model
func (mt *MockTranslator) SupportsModel(model string) bool {
	return strings.EqualFold(mt.name, model) ||
		strings.Contains(strings.ToLower(mt.name), strings.ToLower(model))
}

// mockTranslate generates a mock translation
func (mt *MockTranslator) mockTranslate(text string) string {
	// This is a simple mock translation that just prefixes the text
	// In a real implementation, this would call an actual LLM API
	switch mt.name {
	case "GPT-3.5":
		return fmt.Sprintf("[Translated by GPT-3.5] %s", text)
	case "GPT-4":
		return fmt.Sprintf("[Translated by GPT-4] %s", text)
	case "Claude":
		return fmt.Sprintf("[Translated by Claude] %s", text)
	case "Llama":
		return fmt.Sprintf("[Translated by Llama] %s", text)
	default:
		return fmt.Sprintf("[Translated by %s] %s", mt.name, text)
	}
}
