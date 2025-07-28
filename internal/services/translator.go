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
	// Register mock translators for models that don't have API keys configured
	// In a real implementation, these would be actual LLM API clients when keys are available

	// Register real translators if API keys are configured
	if ts.config.HasOpenAIKey() {
		// TODO: Implement real OpenAI translator when we have the actual implementation
		// For now, we'll still use mock translators but indicate they would be real
		ts.translators["gpt-3.5"] = NewMockTranslator("GPT-3.5 (OpenAI)")
		ts.translators["gpt-4"] = NewMockTranslator("GPT-4 (OpenAI)")
	} else {
		ts.translators["gpt-3.5"] = NewMockTranslator("GPT-3.5")
		ts.translators["gpt-4"] = NewMockTranslator("GPT-4")
	}

	if ts.config.HasAnthropicKey() {
		// TODO: Implement real Anthropic translator when we have the actual implementation
		// For now, we'll still use mock translators but indicate they would be real
		ts.translators["claude"] = NewMockTranslator("Claude (Anthropic)")
	} else {
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
