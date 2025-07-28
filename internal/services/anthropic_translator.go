package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"translator-service/internal/models"
)

// AnthropicTranslator implements the Translator interface for Anthropic models
type AnthropicTranslator struct {
	apiKey   string
	endpoint string
	client   *http.Client
}

// NewAnthropicTranslator creates a new Anthropic translator
func NewAnthropicTranslator(apiKey, endpoint string) *AnthropicTranslator {
	return &AnthropicTranslator{
		apiKey:   apiKey,
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Translate translates text using the Anthropic API
func (at *AnthropicTranslator) Translate(ctx context.Context, req *models.TranslationRequest) (*models.TranslationResponse, error) {
	// Create the Anthropic API request
	apiReq := AnthropicRequest{
		Model:     req.Model,
		Messages:  []AnthropicMessage{{Role: "user", Content: at.createPrompt(req.Text)}},
		MaxTokens: 1000,
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", at.endpoint+"/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", at.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	// Make the API call
	resp, err := at.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make API call: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var apiResp AnthropicResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	// Check for API errors
	if len(apiResp.Error.Type) > 0 {
		return nil, fmt.Errorf("API error (%s): %s", apiResp.Error.Type, apiResp.Error.Message)
	}

	// Extract translation from response
	if len(apiResp.Content) == 0 {
		return nil, fmt.Errorf("API returned no translation content")
	}

	translation := apiResp.Content[0].Text
	if translation == "" {
		return nil, fmt.Errorf("API returned empty translation")
	}

	return &models.TranslationResponse{
		Original:    req.Text,
		Translation: translation,
		Model:       req.Model,
	}, nil
}

// createPrompt creates a prompt for translation
func (at *AnthropicTranslator) createPrompt(text string) string {
	return fmt.Sprintf("Translate the following English text to Chinese. Provide only the translation without any explanation.\n\nEnglish: %s\n\nChinese:", text)
}

// Name returns the name of the translator
func (at *AnthropicTranslator) Name() string {
	return "Anthropic"
}

// SupportsModel returns true if the translator supports the given model
func (at *AnthropicTranslator) SupportsModel(model string) bool {
	switch model {
	case "claude-3-opus-20240229", "claude-3-sonnet-20240229", "claude-3-haiku-20240307", "claude-3-opus", "claude-3-sonnet", "claude-3-haiku":
		return true
	default:
		return false
	}
}

// AnthropicRequest represents the request structure for Anthropic API
type AnthropicRequest struct {
	Model     string             `json:"model"`
	Messages  []AnthropicMessage `json:"messages"`
	MaxTokens int                `json:"max_tokens"`
}

// AnthropicMessage represents a single message in the conversation
type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AnthropicResponse represents the response structure from Anthropic API
type AnthropicResponse struct {
	ID      string             `json:"id"`
	Type    string             `json:"type"`
	Role    string             `json:"role"`
	Content []AnthropicContent `json:"content"`
	Model   string             `json:"model"`
	Error   AnthropicError     `json:"error"`
}

// AnthropicContent represents the content of a message
type AnthropicContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// AnthropicError represents an error returned by the API
type AnthropicError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
