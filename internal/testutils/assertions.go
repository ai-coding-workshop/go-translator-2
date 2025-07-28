package testutils

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

// AssertStatus checks that the response recorder has the expected status code
func AssertStatus(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int) {
	t.Helper()
	if status := rr.Code; status != expectedStatus {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, expectedStatus)
	}
}

// AssertContentType checks that the response has the expected content type
func AssertContentType(t *testing.T, rr *httptest.ResponseRecorder, expectedContentType string) {
	t.Helper()
	contentType := rr.Header().Get("Content-Type")
	if !strings.Contains(contentType, expectedContentType) {
		t.Errorf("Handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}
}

// AssertBodyContains checks that the response body contains the expected substring
func AssertBodyContains(t *testing.T, rr *httptest.ResponseRecorder, expectedSubstring string) {
	t.Helper()
	if !strings.Contains(rr.Body.String(), expectedSubstring) {
		t.Errorf("Handler returned unexpected body: got %v", rr.Body.String())
	}
}

// AssertJSONField checks that the JSON response contains the expected field
func AssertJSONField(t *testing.T, rr *httptest.ResponseRecorder, field string) {
	t.Helper()
	var result map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
		return
	}

	if _, exists := result[field]; !exists {
		t.Errorf("JSON response missing expected field '%s': got %v", field, result)
	}
}

// AssertJSONFieldEquals checks that the JSON response field has the expected value
func AssertJSONFieldEquals(t *testing.T, rr *httptest.ResponseRecorder, field string, expectedValue interface{}) {
	t.Helper()
	var result map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
		return
	}

	if value, exists := result[field]; !exists {
		t.Errorf("JSON response missing expected field '%s': got %v", field, result)
	} else if value != expectedValue {
		t.Errorf("JSON field '%s' has wrong value: got %v want %v", field, value, expectedValue)
	}
}
