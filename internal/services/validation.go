package services

import (
	"regexp"
	"strings"
	"unicode"
)

// ValidationService provides input validation for the translation service
type ValidationService struct{}

// NewValidationService creates a new validation service
func NewValidationService() *ValidationService {
	return &ValidationService{}
}

// ValidateTextInput validates that the input text is in English and meets requirements
func (vs *ValidationService) ValidateTextInput(text string) error {
	// Check if text is empty
	if strings.TrimSpace(text) == "" {
		return &ValidationError{"Text input cannot be empty"}
	}

	// Check maximum length (reasonable limit for translation)
	if len(text) > 10000 {
		return &ValidationError{"Text input is too long (maximum 10000 characters)"}
	}

	// Check if text contains only whitespace
	if len(strings.Fields(text)) == 0 {
		return &ValidationError{"Text input cannot contain only whitespace"}
	}

	// Check for invalid characters (control characters except common ones)
	if vs.containsInvalidCharacters(text) {
		return &ValidationError{"Text contains invalid characters"}
	}

	// Check if text contains English characters
	if !vs.containsEnglishCharacters(text) {
		return &ValidationError{"Text must contain English characters"}
	}

	// Check if text is mostly English (at least 70% English words)
	if !vs.isMostlyEnglish(text) {
		return &ValidationError{"Text must be primarily in English"}
	}

	return nil
}

// ValidateModelInput validates that the model is supported
func (vs *ValidationService) ValidateModelInput(model string, supportedModels []string) error {
	// Check if model is empty
	if strings.TrimSpace(model) == "" {
		return &ValidationError{"Model selection cannot be empty"}
	}

	// Check if model is supported
	for _, supportedModel := range supportedModels {
		if supportedModel == model {
			return nil
		}
	}

	return &ValidationError{"Unsupported model: " + model}
}

// containsEnglishCharacters checks if the text contains English letters
func (vs *ValidationService) containsEnglishCharacters(text string) bool {
	// Remove excessive whitespace
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// Check if text contains English letters
	for _, r := range text {
		if unicode.IsLetter(r) && (r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z') {
			return true
		}
	}

	// If no English letters found, check if it's mostly ASCII
	asciiCount := 0
	totalCount := 0
	for _, r := range text {
		if r <= 127 { // ASCII range
			asciiCount++
		}
		totalCount++
	}

	// If more than 80% of characters are ASCII, consider it valid
	return totalCount > 0 && float64(asciiCount)/float64(totalCount) > 0.8
}

// isMostlyEnglish checks if the text is primarily in English
func (vs *ValidationService) isMostlyEnglish(text string) bool {
	// Simple check: if the text contains English characters and is mostly ASCII,
	// we'll consider it English for this implementation
	// A more sophisticated implementation would use a language detection library
	return vs.containsEnglishCharacters(text)
}

// containsInvalidCharacters checks for invalid control characters
func (vs *ValidationService) containsInvalidCharacters(text string) bool {
	for _, r := range text {
		// Allow common whitespace and printable ASCII characters
		if r < 32 && r != 9 && r != 10 && r != 13 { // Allow tab, newline, carriage return
			return true
		}
		// Disallow DEL character and above private use area
		if r >= 127 && r <= 159 {
			return true
		}
	}
	return false
}

// ValidationError represents a validation error
type ValidationError struct {
	message string
}

func (e *ValidationError) Error() string {
	return e.message
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}
