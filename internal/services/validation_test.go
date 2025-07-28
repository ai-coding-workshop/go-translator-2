package services

import (
	"strings"
	"testing"
)

func TestValidationService_ValidateTextInput(t *testing.T) {
	vs := NewValidationService()

	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "Valid English text",
			input:       "Hello, world!",
			expectError: false,
		},
		{
			name:        "Empty text",
			input:       "",
			expectError: true,
		},
		{
			name:        "Whitespace only",
			input:       "   \t\n  ",
			expectError: true,
		},
		{
			name:        "Text too long",
			input:       string(make([]byte, 10001)),
			expectError: true,
		},
		{
			name:        "Valid length text",
			input:       "Hello, this is a valid English text that is long enough to test the length validation but not too long to exceed the limit. " + strings.Repeat("This is repeated text. ", 100),
			expectError: false,
		},
		{
			name:        "Text with invalid control characters",
			input:       "Hello\x01World",
			expectError: true,
		},
		{
			name:        "Text without English characters",
			input:       "你好世界",
			expectError: true,
		},
		{
			name:        "Mixed text with mostly English",
			input:       "Hello 你好 world",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := vs.ValidateTextInput(tt.input)
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidationService_ValidateModelInput(t *testing.T) {
	vs := NewValidationService()
	supportedModels := []string{"gpt-3.5", "gpt-4", "claude", "llama"}

	tests := []struct {
		name            string
		model           string
		supportedModels []string
		expectError     bool
	}{
		{
			name:            "Valid model",
			model:           "gpt-3.5",
			supportedModels: supportedModels,
			expectError:     false,
		},
		{
			name:            "Empty model",
			model:           "",
			supportedModels: supportedModels,
			expectError:     true,
		},
		{
			name:            "Unsupported model",
			model:           "unsupported-model",
			supportedModels: supportedModels,
			expectError:     true,
		},
		{
			name:            "Case sensitive model",
			model:           "GPT-3.5",
			supportedModels: supportedModels,
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := vs.ValidateModelInput(tt.model, tt.supportedModels)
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestIsValidationError(t *testing.T) {
	vs := NewValidationService()

	// Test with a validation error
	err := vs.ValidateTextInput("")
	if !IsValidationError(err) {
		t.Errorf("Expected validation error to be detected as validation error")
	}

	// Test with a non-validation error
	if IsValidationError(nil) {
		t.Errorf("Expected nil to not be a validation error")
	}
}
