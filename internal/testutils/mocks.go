package testutils

import (
	"context"
	"time"

	"translator-service/internal/models"
)

// MockTranslator is a mock implementation of the Translator interface for testing
type MockTranslator struct {
	NameFunc           func() string
	SupportsModelFunc  func(model string) bool
	TranslateFunc      func(ctx context.Context, req *models.TranslationRequest) (*models.TranslationResponse, error)
	NameValue          string
	SupportsModelValue bool
	TranslateResponse  *models.TranslationResponse
	TranslateError     error
}

// Name returns the name of the translator
func (m *MockTranslator) Name() string {
	if m.NameFunc != nil {
		return m.NameFunc()
	}
	return m.NameValue
}

// SupportsModel returns true if the translator supports the given model
func (m *MockTranslator) SupportsModel(model string) bool {
	if m.SupportsModelFunc != nil {
		return m.SupportsModelFunc(model)
	}
	return m.SupportsModelValue
}

// Translate performs a translation
func (m *MockTranslator) Translate(ctx context.Context, req *models.TranslationRequest) (*models.TranslationResponse, error) {
	if m.TranslateFunc != nil {
		return m.TranslateFunc(ctx, req)
	}
	return m.TranslateResponse, m.TranslateError
}

// MockContextWithTimeout creates a context with a timeout for testing
func MockContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// MockTranslationRequest creates a translation request for testing
func MockTranslationRequest(text, model string) *models.TranslationRequest {
	return &models.TranslationRequest{
		Text:  text,
		Model: model,
	}
}

// MockTranslationResponse creates a translation response for testing
func MockTranslationResponse(original, translation, model string) *models.TranslationResponse {
	return &models.TranslationResponse{
		Original:    original,
		Translation: translation,
		Model:       model,
	}
}
