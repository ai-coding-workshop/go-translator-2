package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"translator-service/internal/models"
	"translator-service/internal/services"
)

var (
	homeTemplate   *template.Template
	resultTemplate *template.Template
)

func init() {
	// Parse templates
	homeTemplate = template.Must(template.ParseFiles(filepath.Join("web", "templates", "home.html")))
	resultTemplate = template.Must(template.ParseFiles(filepath.Join("web", "templates", "result.html")))
}

// HomeHandler serves the main web page
func NewHomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// Render home template
		if err := homeTemplate.Execute(w, nil); err != nil {
			log.Printf("Error rendering home template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

// TranslateHandler processes translation requests from the web form
type TranslateHandler struct {
	translatorService *services.TranslatorService
}

func NewTranslateHandler(translatorService *services.TranslatorService) http.HandlerFunc {
	handler := &TranslateHandler{
		translatorService: translatorService,
	}

	return handler.ServeHTTP
}

func (h *TranslateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get form values
	text := strings.TrimSpace(r.FormValue("text"))
	model := strings.TrimSpace(r.FormValue("model"))

	// Validate input
	if text == "" {
		http.Error(w, "Please enter text to translate", http.StatusBadRequest)
		return
	}

	if model == "" {
		http.Error(w, "Please select a translation model", http.StatusBadRequest)
		return
	}

	// Create translation request
	req := &models.TranslationRequest{
		Text:  text,
		Model: model,
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Perform translation
	response, err := h.translatorService.Translate(ctx, req)
	if err != nil {
		log.Printf("Translation error: %v", err)
		// Provide user-friendly error message
		if strings.Contains(err.Error(), "validation error") {
			http.Error(w, fmt.Sprintf("Invalid input: %s", strings.TrimPrefix(err.Error(), "validation error: ")), http.StatusBadRequest)
		} else if strings.Contains(err.Error(), "unsupported model") {
			http.Error(w, "Selected translation model is not supported", http.StatusBadRequest)
		} else if strings.Contains(err.Error(), "context deadline exceeded") {
			http.Error(w, "Translation request timed out. Please try again.", http.StatusRequestTimeout)
		} else if strings.Contains(err.Error(), "context canceled") {
			http.Error(w, "Translation request was canceled", http.StatusBadRequest)
		} else {
			http.Error(w, "Translation service temporarily unavailable. Please try again in a moment.", http.StatusServiceUnavailable)
		}
		return
	}

	// Render result template
	if err := resultTemplate.Execute(w, response); err != nil {
		log.Printf("Error rendering result template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// APIHandler handles REST API requests for translation
type APIHandler struct {
	translatorService *services.TranslatorService
}

func NewAPIHandler(translatorService *services.TranslatorService) http.HandlerFunc {
	handler := &APIHandler{
		translatorService: translatorService,
	}

	return handler.ServeHTTP
}

func (h *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode JSON request
	var req models.TranslationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	// Trim whitespace
	req.Text = strings.TrimSpace(req.Text)
	req.Model = strings.TrimSpace(req.Model)

	// Validate request
	if req.Text == "" {
		http.Error(w, "Text field is required", http.StatusBadRequest)
		return
	}

	if req.Model == "" {
		http.Error(w, "Model field is required", http.StatusBadRequest)
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Perform translation
	response, err := h.translatorService.Translate(ctx, &req)
	if err != nil {
		log.Printf("Translation error: %v", err)
		// Provide structured error response
		errorResponse := map[string]interface{}{
			"error":   true,
			"message": h.getErrorMessage(err),
			"details": err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(h.getErrorCode(err))
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			log.Printf("Error encoding JSON error response: %v", err)
		}
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// getErrorMessage returns a user-friendly error message based on the error
func (h *APIHandler) getErrorMessage(err error) string {
	if strings.Contains(err.Error(), "validation error") {
		return fmt.Sprintf("Invalid input: %s", strings.TrimPrefix(err.Error(), "validation error: "))
	} else if strings.Contains(err.Error(), "unsupported model") {
		return "Selected translation model is not supported"
	} else if strings.Contains(err.Error(), "context deadline exceeded") {
		return "Translation request timed out"
	} else if strings.Contains(err.Error(), "context canceled") {
		return "Translation request was canceled"
	} else {
		return "Translation service temporarily unavailable"
	}
}

// getErrorCode returns an appropriate HTTP status code based on the error
func (h *APIHandler) getErrorCode(err error) int {
	if strings.Contains(err.Error(), "validation error") {
		return http.StatusBadRequest
	} else if strings.Contains(err.Error(), "unsupported model") {
		return http.StatusBadRequest
	} else if strings.Contains(err.Error(), "context deadline exceeded") {
		return http.StatusRequestTimeout
	} else if strings.Contains(err.Error(), "context canceled") {
		return http.StatusBadRequest
	} else {
		return http.StatusServiceUnavailable
	}
}
