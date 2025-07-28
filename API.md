# API Documentation

This document describes the REST API endpoints available in the Translation Service.

## Table of Contents
- [Base URL](#base-url)
- [Endpoints](#endpoints)
  - [Web Interface](#web-interface)
  - [Translation API](#translation-api)
- [Request/Response Formats](#requestresponse-formats)
- [Error Handling](#error-handling)
- [Examples](#examples)

## Base URL

All API endpoints are relative to the base URL where the service is hosted, typically:
```
http://localhost:8080
```

## Endpoints

### Web Interface

#### GET /
Serves the main web page with the translation form.

**Response:**
- 200 OK - HTML page with translation form
- 500 Internal Server Error - If there's an issue rendering the template

#### POST /translate
Processes translation requests from the web form.

**Request Parameters:**
- `text` (string, required) - The English text to translate
- `model` (string, required) - The LLM model to use for translation

**Response:**
- 200 OK - HTML page with translation results
- 400 Bad Request - If text or model parameters are missing or invalid
- 405 Method Not Allowed - If method is not POST
- 500 Internal Server Error - If there's an issue with the translation service

### Translation API

#### POST /api/translate
REST API endpoint for translating text programmatically.

**Request Format:**
```json
{
  "text": "string",
  "model": "string"
}
```

**Request Fields:**
- `text` (string, required) - The English text to translate
- `model` (string, required) - The LLM model to use for translation

**Supported Models:**
- `gpt-4` - OpenAI GPT-4
- `gpt-3.5` - OpenAI GPT-3.5
- `claude-3-opus` - Anthropic Claude 3 Opus
- `claude-3-sonnet` - Anthropic Claude 3 Sonnet
- `claude-3-haiku` - Anthropic Claude 3 Haiku

**Response Format (Success):**
```json
{
  "original": "string",
  "translation": "string",
  "model": "string"
}
```

**Response Fields:**
- `original` - The original text that was translated
- `translation` - The translated text
- `model` - The model that was used for translation

**Response Format (Error):**
```json
{
  "error": true,
  "message": "string",
  "details": "string"
}
```

**HTTP Status Codes:**
- 200 OK - Translation successful
- 400 Bad Request - Invalid request data
- 405 Method Not Allowed - Wrong HTTP method
- 408 Request Timeout - Translation request timed out
- 503 Service Unavailable - Translation service temporarily unavailable
- 500 Internal Server Error - Unexpected server error

## Request/Response Formats

All API requests and responses use JSON format with UTF-8 encoding.

**Content-Type Header:**
```
Content-Type: application/json
```

## Error Handling

The API uses standard HTTP status codes to indicate the success or failure of requests.

**Common Error Responses:**

400 Bad Request:
```json
{
  "error": true,
  "message": "Text field is required",
  "details": "Text field is required"
}
```

408 Request Timeout:
```json
{
  "error": true,
  "message": "Translation request timed out",
  "details": "context deadline exceeded"
}
```

503 Service Unavailable:
```json
{
  "error": true,
  "message": "Translation service temporarily unavailable",
  "details": "API error details"
}
```

## Examples

### Web Form Translation

```bash
curl -X POST http://localhost:8080/translate \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "text=Hello%2C%20world%21&model=gpt-3.5"
```

### REST API Translation

```bash
curl -X POST http://localhost:8080/api/translate \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Hello, world!",
    "model": "gpt-3.5"
  }'
```

**Response:**
```json
{
  "original": "Hello, world!",
  "translation": "你好，世界！",
  "model": "gpt-3.5"
}
```

### Error Response Example

```bash
curl -X POST http://localhost:8080/api/translate \
  -H "Content-Type: application/json" \
  -d '{
    "text": "",
    "model": "gpt-3.5"
  }'
```

**Response:**
```json
{
  "error": true,
  "message": "Invalid input: Text field is required",
  "details": "validation error: Text field is required"
}
```
