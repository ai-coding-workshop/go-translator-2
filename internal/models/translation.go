package models

import "context"

// TranslationRequest represents a request to translate text
type TranslationRequest struct {
	Text  string `json:"text"`
	Model string `json:"model"`
}

// TranslationResponse represents a translation response
type TranslationResponse struct {
	Original    string `json:"original"`
	Translation string `json:"translation"`
	Model       string `json:"model"`
}

// Translator defines the interface for translation services
type Translator interface {
	// Translate translates the given text using the specified model
	Translate(ctx context.Context, req *TranslationRequest) (*TranslationResponse, error)

	// Name returns the name of the translator
	Name() string

	// SupportsModel returns true if the translator supports the given model
	SupportsModel(model string) bool
}
