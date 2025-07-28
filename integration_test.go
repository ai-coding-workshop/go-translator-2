package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"translator-service/internal/config"
	"translator-service/internal/handlers"
	"translator-service/internal/models"
	"translator-service/internal/services"
)

func TestFullServiceIntegration(t *testing.T) {
	// Set test mode environment variable
	os.Setenv("GO_TEST_MODE", "true")
	defer os.Unsetenv("GO_TEST_MODE")

	// Create a test config
	cfg := &config.Config{
		ServerPort: "8080",
		Timeout:    30,
	}

	// Create translator service
	translatorService := services.NewTranslatorService(cfg)

	// Create handlers
	homeHandler := handlers.NewHomeHandler()
	translateHandler := handlers.NewTranslateHandler(translatorService)
	apiHandler := handlers.NewAPIHandler(translatorService)

	// Create a test server with our handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/translate", translateHandler)
	mux.HandleFunc("/api/translate", apiHandler)

	// Test the home page
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	homeHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("HomeHandler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Test the API endpoint with a valid request
	requestData := models.TranslationRequest{
		Text:  "Hello, world!",
		Model: "gpt-3.5",
	}
	jsonData, _ := json.Marshal(requestData)

	req, err = http.NewRequest("POST", "/api/translate", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	apiHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("APIHandler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Parse the response
	var translationResponse models.TranslationResponse
	if err := json.NewDecoder(rr.Body).Decode(&translationResponse); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if translationResponse.Original != "Hello, world!" {
		t.Errorf("Expected original text to match, got %s", translationResponse.Original)
	}

	// Test the API endpoint with invalid JSON
	req, err = http.NewRequest("POST", "/api/translate", bytes.NewBufferString("{invalid json}"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	apiHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("APIHandler returned wrong status code for invalid JSON: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Test the API endpoint with empty text
	requestData = models.TranslationRequest{
		Text:  "",
		Model: "gpt-3.5",
	}
	jsonData, _ = json.Marshal(requestData)

	req, err = http.NewRequest("POST", "/api/translate", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	apiHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("APIHandler returned wrong status code for empty text: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Test the translate form endpoint with empty text
	form := strings.NewReader("text=&model=gpt-3.5")
	req, err = http.NewRequest("POST", "/translate", form)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	translateHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("TranslateHandler returned wrong status code for empty text: got %v want %v",
			status, http.StatusBadRequest)
	}
}
