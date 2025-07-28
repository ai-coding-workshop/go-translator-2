package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

// CreateJSONRequest creates an HTTP request with JSON body for testing
func CreateJSONRequest(method, url string, data interface{}) (*http.Request, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// CreateFormRequest creates an HTTP request with form data for testing
func CreateFormRequest(method, url, formData string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(formData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

// CreateFormRequestFromValues creates an HTTP request with form data from url.Values for testing
func CreateFormRequestFromValues(method, url string, values url.Values) (*http.Request, error) {
	return CreateFormRequest(method, url, values.Encode())
}

// ExecuteHandlerAndRecord executes an HTTP handler and returns the response recorder
func ExecuteHandlerAndRecord(handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// ExecuteHandlerAndReturnBody executes an HTTP handler and returns the response body as string
func ExecuteHandlerAndReturnBody(handler http.HandlerFunc, req *http.Request) string {
	rr := ExecuteHandlerAndRecord(handler, req)
	return rr.Body.String()
}

// ExecuteHandlerAndReturnStatus executes an HTTP handler and returns the status code
func ExecuteHandlerAndReturnStatus(handler http.HandlerFunc, req *http.Request) int {
	rr := ExecuteHandlerAndRecord(handler, req)
	return rr.Code
}
