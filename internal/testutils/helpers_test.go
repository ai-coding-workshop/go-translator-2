package testutils

import (
	"net/url"
	"testing"
)

func TestCreateFormRequestFromValues(t *testing.T) {
	values := url.Values{}
	values.Set("text", "Hello, world!")
	values.Set("model", "gpt-3.5")

	req, err := CreateFormRequestFromValues("POST", "/translate", values)
	if err != nil {
		t.Fatalf("Failed to create form request: %v", err)
	}

	if req.Method != "POST" {
		t.Errorf("Expected method POST, got %s", req.Method)
	}

	if req.URL.Path != "/translate" {
		t.Errorf("Expected path /translate, got %s", req.URL.Path)
	}

	if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Errorf("Expected content type application/x-www-form-urlencoded, got %s", req.Header.Get("Content-Type"))
	}
}
