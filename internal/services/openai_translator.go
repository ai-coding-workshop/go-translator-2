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

// OpenAITranslator implements the Translator interface for OpenAI models
type OpenAITranslator struct {
	apiKey   string
	endpoint string
	client   *http.Client
}

// NewOpenAITranslator creates a new OpenAI translator
func NewOpenAITranslator(apiKey, endpoint string) *OpenAITranslator {
	return &OpenAITranslator{
		apiKey:   apiKey,
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Translate translates text using the OpenAI API
func (ot *OpenAITranslator) Translate(ctx context.Context, req *models.TranslationRequest) (*models.TranslationResponse, error) {
	// Create the OpenAI API request
	apiReq := OpenAIRequest{
		Model: req.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a professional English to Chinese translator. Translate the following English text to Chinese. Provide only the translation without any explanation.",
			},
			{
				Role:    "user",
				Content: req.Text,
			},
		},
		Temperature: 0.3,
		MaxTokens:   1000,
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", ot.endpoint+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+ot.apiKey)

	// Make the API call
	resp, err := ot.client.Do(httpReq)
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
	var apiResp OpenAIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	// Check for API errors
	if len(apiResp.Error.Message) > 0 {
		return nil, fmt.Errorf("API error: %s", apiResp.Error.Message)
	}

	// Extract translation from response
	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("API returned no translation choices")
	}

	translation := apiResp.Choices[0].Message.Content
	if translation == "" {
		return nil, fmt.Errorf("API returned empty translation")
	}

	return &models.TranslationResponse{
		Original:    req.Text,
		Translation: translation,
		Model:       req.Model,
	}, nil
}

// Name returns the name of the translator
func (ot *OpenAITranslator) Name() string {
	return "OpenAI"
}

// SupportsModel returns true if the translator supports the given model
func (ot *OpenAITranslator) SupportsModel(model string) bool {
	switch model {
	case "gpt-3.5-turbo", "gpt-3.5", "gpt-4", "gpt-4-turbo", "gpt-4o":
		return true
	default:
		return false
	}
}

// OpenAIRequest represents the request structure for OpenAI API
type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// Message represents a single message in the conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents the response structure from OpenAI API
type OpenAIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
	Error   APIError `json:"error"`
}

// Choice represents a single choice in the API response
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// APIError represents an error returned by the API
type APIError struct {
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Param   interface{} `json:"param"`
	Code    interface{} `json:"code"`
}
