package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"translator-service/internal/config"
	"translator-service/internal/models"
)

// TranslatorService manages multiple translation providers
type TranslatorService struct {
	translators       map[string]models.Translator
	validationService *ValidationService
	config            *config.Config
}

// NewTranslatorService creates a new translator service
func NewTranslatorService(cfg *config.Config) *TranslatorService {
	service := &TranslatorService{
		translators:       make(map[string]models.Translator),
		validationService: NewValidationService(),
		config:            cfg,
	}

	// Register supported translators
	service.registerTranslators()

	return service
}

// registerTranslators registers all supported translation providers
func (ts *TranslatorService) registerTranslators() {
	// Register real translators if API keys are configured
	if ts.config.HasOpenAIKey() {
		openaiTranslator := NewOpenAITranslator(ts.config.GetOpenAIKey(), ts.config.GetOpenAIEndpoint())

		// Support for Qwen models (Alibaba Cloud) - using OpenAI-compatible API
		if strings.HasPrefix(ts.config.GetOpenAIEndpoint(), "https://idealab.alibaba-inc.com") {
			ts.translators["Qwen3-Coder-Plus"] = openaiTranslator
			ts.translators["qwen-max-latest"] = openaiTranslator
			ts.translators["qwen-plus"] = openaiTranslator
			ts.translators["qwen2.5-max"] = openaiTranslator
			ts.translators["qwen2.5-plus"] = openaiTranslator
		} else {
			ts.translators["gpt-3.5-turbo"] = openaiTranslator
			ts.translators["gpt-3.5"] = openaiTranslator
			ts.translators["gpt-4"] = openaiTranslator
			ts.translators["gpt-4-turbo"] = openaiTranslator
			ts.translators["gpt-4o"] = openaiTranslator
		}
	} else {
		// Register mock translators for OpenAI models when no API key is present
		ts.translators["gpt-3.5-turbo"] = NewMockTranslator("GPT-3.5 Turbo")
		ts.translators["gpt-3.5"] = NewMockTranslator("GPT-3.5")
		ts.translators["gpt-4"] = NewMockTranslator("GPT-4")
		ts.translators["gpt-4-turbo"] = NewMockTranslator("GPT-4 Turbo")
		ts.translators["gpt-4o"] = NewMockTranslator("GPT-4O")
		// Also register mock translators for Qwen models
		ts.translators["Qwen3-Coder-Plus"] = NewMockTranslator("Qwen3 Coder Plus")
		ts.translators["qwen-max-latest"] = NewMockTranslator("Qwen Max Latest")
		ts.translators["qwen-plus"] = NewMockTranslator("Qwen Plus")
		ts.translators["qwen2.5-max"] = NewMockTranslator("Qwen 2.5 Max")
		ts.translators["qwen2.5-plus"] = NewMockTranslator("Qwen 2.5 Plus")
	}

	if ts.config.HasAnthropicKey() {
		anthropicTranslator := NewAnthropicTranslator(ts.config.GetAnthropicKey(), ts.config.GetAnthropicEndpoint())
		ts.translators["claude-3-opus"] = anthropicTranslator
		ts.translators["claude-3-sonnet"] = anthropicTranslator
		ts.translators["claude-3-haiku"] = anthropicTranslator
		ts.translators["claude-3-opus-20240229"] = anthropicTranslator
		ts.translators["claude-3-sonnet-20240229"] = anthropicTranslator
		ts.translators["claude-3-haiku-20240307"] = anthropicTranslator
		ts.translators["claude"] = anthropicTranslator
	} else {
		// Register mock translators for Anthropic models when no API key is present
		ts.translators["claude-3-opus"] = NewMockTranslator("Claude 3 Opus")
		ts.translators["claude-3-sonnet"] = NewMockTranslator("Claude 3 Sonnet")
		ts.translators["claude-3-haiku"] = NewMockTranslator("Claude 3 Haiku")
		ts.translators["claude-3-opus-20240229"] = NewMockTranslator("Claude 3 Opus (2024-02-29)")
		ts.translators["claude-3-sonnet-20240229"] = NewMockTranslator("Claude 3 Sonnet (2024-02-29)")
		ts.translators["claude-3-haiku-20240307"] = NewMockTranslator("Claude 3 Haiku (2024-03-07)")
		ts.translators["claude"] = NewMockTranslator("Claude")
	}

	// Llama is always a mock translator (open source model)
	ts.translators["llama"] = NewMockTranslator("Llama")
}

// Translate translates text using the specified model with retry logic
func (ts *TranslatorService) Translate(ctx context.Context, req *models.TranslationRequest) (*models.TranslationResponse, error) {
	// Validate input
	if err := ts.validationService.ValidateTextInput(req.Text); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if err := ts.validationService.ValidateModelInput(req.Model, ts.GetSupportedModels()); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Find the appropriate translator
	translator, exists := ts.translators[req.Model]
	if !exists {
		return nil, fmt.Errorf("unsupported model: %s", req.Model)
	}

	// Perform translation with retry logic
	var response *models.TranslationResponse
	var err error

	// Retry up to 3 times for transient errors
	for attempt := 0; attempt < 3; attempt++ {
		response, err = translator.Translate(ctx, req)
		if err == nil {
			// Success
			return response, nil
		}

		// Log the error
		log.Printf("Translation attempt %d failed with %s: %v", attempt+1, req.Model, err)

		// Don't retry on validation errors or context cancellation
		if ctx.Err() != nil {
			break
		}

		// Check if it's a validation error by checking the error message
		if err != nil && strings.Contains(err.Error(), "validation error:") {
			break
		}

		// For transient network errors, retry after a short delay
		if attempt < 2 { // Don't sleep on the last attempt
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt+1) * 100 * time.Millisecond): // Exponential backoff
				// Continue to retry
			}
		}
	}

	if err != nil {
		log.Printf("Translation failed after retries with %s: %v", req.Model, err)
		return nil, fmt.Errorf("failed to translate with %s after retries: %w", req.Model, err)
	}

	return response, nil
}

// GetSupportedModels returns a list of supported models
func (ts *TranslatorService) GetSupportedModels() []string {
	models := make([]string, 0, len(ts.translators))
	for model := range ts.translators {
		models = append(models, model)
	}
	return models
}

// IsModelSupported checks if a model is supported
func (ts *TranslatorService) IsModelSupported(model string) bool {
	_, exists := ts.translators[model]
	return exists
}
