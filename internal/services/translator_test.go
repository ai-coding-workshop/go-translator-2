package services

import (
	"context"
	"testing"
	"time"

	"translator-service/internal/config"
	"translator-service/internal/models"
)

// MockTranslatorForTesting is a mock translator for testing purposes
type MockTranslatorForTesting struct {
	name          string
	translateFunc func(ctx context.Context, req *models.TranslationRequest) (*models.TranslationResponse, error)
}

func (m *MockTranslatorForTesting) Translate(ctx context.Context, req *models.TranslationRequest) (*models.TranslationResponse, error) {
	if m.translateFunc != nil {
		return m.translateFunc(ctx, req)
	}
	return &models.TranslationResponse{
		Original:    req.Text,
		Translation: "mock translation",
		Model:       m.name,
	}, nil
}

func (m *MockTranslatorForTesting) Name() string {
	return m.name
}

func (m *MockTranslatorForTesting) SupportsModel(model string) bool {
	return m.name == model
}

func TestTranslatorService_Translate(t *testing.T) {
	// Create a test config
	cfg := &config.Config{
		ServerPort: "8080",
		Timeout:    30,
	}

	// Create translator service
	ts := NewTranslatorService(cfg)

	// Register a mock translator for testing
	mockTranslator := &MockTranslatorForTesting{name: "test-model"}
	ts.translators["test-model"] = mockTranslator

	tests := []struct {
		name          string
		request       *models.TranslationRequest
		setupMock     func()
		expectError   bool
		expectTimeout bool
	}{
		{
			name: "Valid translation",
			request: &models.TranslationRequest{
				Text:  "Hello, world!",
				Model: "test-model",
			},
			expectError: false,
		},
		{
			name: "Empty text",
			request: &models.TranslationRequest{
				Text:  "",
				Model: "test-model",
			},
			expectError: true,
		},
		{
			name: "Unsupported model",
			request: &models.TranslationRequest{
				Text:  "Hello, world!",
				Model: "unsupported-model",
			},
			expectError: true,
		},
		{
			name: "Context timeout",
			request: &models.TranslationRequest{
				Text:  "Hello, world!",
				Model: "test-model",
			},
			setupMock: func() {
				mockTranslator.translateFunc = func(ctx context.Context, req *models.TranslationRequest) (*models.TranslationResponse, error) {
					// Simulate a long operation
					time.Sleep(100 * time.Millisecond)
					return nil, context.DeadlineExceeded
				}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
				defer func() {
					mockTranslator.translateFunc = nil
				}()
			}

			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()

			_, err := ts.Translate(ctx, tt.request)
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestTranslatorService_GetSupportedModels(t *testing.T) {
	cfg := &config.Config{
		ServerPort: "8080",
		Timeout:    30,
	}

	ts := NewTranslatorService(cfg)

	// Check that we have the expected number of models (from the automatic registration)
	models := ts.GetSupportedModels()
	if len(models) < 4 {
		t.Errorf("Expected at least 4 supported models, got %d", len(models))
	}

	// Check that some expected models are present
	foundGPT35 := false
	foundGPT4 := false
	foundClaude := false
	foundLlama := false
	for _, model := range models {
		if model == "gpt-3.5" {
			foundGPT35 = true
		}
		if model == "gpt-4" {
			foundGPT4 = true
		}
		if model == "claude" {
			foundClaude = true
		}
		if model == "llama" {
			foundLlama = true
		}
	}

	if !foundGPT35 || !foundGPT4 || !foundClaude || !foundLlama {
		t.Errorf("Expected to find gpt-3.5, gpt-4, claude, and llama in supported models")
	}
}

func TestTranslatorService_IsModelSupported(t *testing.T) {
	cfg := &config.Config{
		ServerPort: "8080",
		Timeout:    30,
	}

	ts := NewTranslatorService(cfg)

	// Add a mock translator
	ts.translators["test-model"] = &MockTranslatorForTesting{name: "test-model"}

	if !ts.IsModelSupported("test-model") {
		t.Errorf("Expected test-model to be supported")
	}

	if ts.IsModelSupported("unsupported-model") {
		t.Errorf("Expected unsupported-model to not be supported")
	}
}

func TestTranslatorService_RetryLogic(t *testing.T) {
	cfg := &config.Config{
		ServerPort: "8080",
		Timeout:    30,
	}

	ts := NewTranslatorService(cfg)

	// Add a mock translator that fails twice then succeeds
	callCount := 0
	mockTranslator := &MockTranslatorForTesting{
		name: "retry-model",
		translateFunc: func(ctx context.Context, req *models.TranslationRequest) (*models.TranslationResponse, error) {
			callCount++
			if callCount < 3 {
				return nil, &ValidationError{"temporary error"}
			}
			return &models.TranslationResponse{
				Original:    req.Text,
				Translation: "success after retries",
				Model:       "retry-model",
			}, nil
		},
	}
	ts.translators["retry-model"] = mockTranslator

	request := &models.TranslationRequest{
		Text:  "Hello, world!",
		Model: "retry-model",
	}

	ctx := context.Background()
	response, err := ts.Translate(ctx, request)

	if err != nil {
		t.Errorf("Expected success after retries, got error: %v", err)
	}

	if response.Translation != "success after retries" {
		t.Errorf("Expected successful translation after retries")
	}

	if callCount != 3 {
		t.Errorf("Expected 3 calls, got %d", callCount)
	}
}
