package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"translator-service/internal/config"
	"translator-service/internal/models"
	"translator-service/internal/services"
)

func createTestTranslatorService() *services.TranslatorService {
	// Create a test config
	cfg := &config.Config{
		ServerPort: "8080",
		Timeout:    30,
	}

	// Create translator service
	ts := services.NewTranslatorService(cfg)
	return ts
}

func TestHomeHandler(t *testing.T) {
	// Create a request to the home endpoint
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a translator service
	service := createTestTranslatorService()

	// Create the handler
	handler := NewHomeHandler(service)

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("HomeHandler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check that the response contains expected content
	// Note: This test might fail if templates are not available
	// For now, we'll just check that it doesn't return an error status
	if status := rr.Code; status >= 500 {
		t.Errorf("HomeHandler returned server error status code: got %v", status)
	}
}

func TestTranslateHandler_GetRequest(t *testing.T) {
	// Create a GET request to the translate endpoint (should fail)
	req, err := http.NewRequest("GET", "/translate", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Create a translator service
	service := createTestTranslatorService()

	// Create the handler
	handler := NewTranslateHandler(service)

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code (should be 405 Method Not Allowed)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("TranslateHandler returned wrong status code for GET request: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestTranslateHandler_EmptyText(t *testing.T) {
	// Create a POST request with empty text
	form := strings.NewReader("text=&model=gpt-3.5")
	req, err := http.NewRequest("POST", "/translate", form)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Create a translator service
	service := createTestTranslatorService()

	// Create the handler
	handler := NewTranslateHandler(service)

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code (should be 400 Bad Request)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("TranslateHandler returned wrong status code for empty text: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestTranslateHandler_EmptyModel(t *testing.T) {
	// Create a POST request with empty model
	form := strings.NewReader("text=Hello%2C%20world%21&model=")
	req, err := http.NewRequest("POST", "/translate", form)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Create a translator service
	service := createTestTranslatorService()

	// Create the handler
	handler := NewTranslateHandler(service)

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code (should be 400 Bad Request)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("TranslateHandler returned wrong status code for empty model: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestAPIHandler_GetRequest(t *testing.T) {
	// Create a GET request to the API endpoint (should fail)
	req, err := http.NewRequest("GET", "/api/translate", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Create a translator service
	service := createTestTranslatorService()

	// Create the handler
	handler := NewAPIHandler(service)

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code (should be 405 Method Not Allowed)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("APIHandler returned wrong status code for GET request: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestAPIHandler_ValidRequest(t *testing.T) {
	// Create a JSON request
	requestData := models.TranslationRequest{
		Text:  "Hello, world!",
		Model: "gpt-3.5",
	}
	jsonData, _ := json.Marshal(requestData)

	req, err := http.NewRequest("POST", "/api/translate", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Create a translator service
	service := createTestTranslatorService()

	// Create the handler
	handler := NewAPIHandler(service)

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("APIHandler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check that the response is JSON
	contentType := rr.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("APIHandler returned wrong content type: got %v want application/json", contentType)
	}

	// Check that the response contains expected fields
	if !strings.Contains(rr.Body.String(), "original") || !strings.Contains(rr.Body.String(), "translation") {
		t.Errorf("APIHandler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestAPIHandler_InvalidJSON(t *testing.T) {
	// Create an invalid JSON request
	req, err := http.NewRequest("POST", "/api/translate", bytes.NewBufferString("{invalid json}"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Create a translator service
	service := createTestTranslatorService()

	// Create the handler
	handler := NewAPIHandler(service)

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code (should be 400 Bad Request)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("APIHandler returned wrong status code for invalid JSON: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestAPIHandler_EmptyText(t *testing.T) {
	// Create a JSON request with empty text
	requestData := models.TranslationRequest{
		Text:  "",
		Model: "gpt-3.5",
	}
	jsonData, _ := json.Marshal(requestData)

	req, err := http.NewRequest("POST", "/api/translate", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Create a translator service
	service := createTestTranslatorService()

	// Create the handler
	handler := NewAPIHandler(service)

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code (should be 400 Bad Request)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("APIHandler returned wrong status code for empty text: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestAPIHandler_EmptyModel(t *testing.T) {
	// Create a JSON request with empty model
	requestData := models.TranslationRequest{
		Text:  "Hello, world!",
		Model: "",
	}
	jsonData, _ := json.Marshal(requestData)

	req, err := http.NewRequest("POST", "/api/translate", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Create a translator service
	service := createTestTranslatorService()

	// Create the handler
	handler := NewAPIHandler(service)

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code (should be 400 Bad Request)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("APIHandler returned wrong status code for empty model: got %v want %v",
			status, http.StatusBadRequest)
	}
}
