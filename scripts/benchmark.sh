#!/bin/bash

# Benchmark script for the Translation Service

echo "Translation Service Benchmark"
echo "============================="

# Check if wrk is installed
if ! command -v wrk &> /dev/null
then
    echo "wrk could not be found. Please install it to run benchmarks."
    echo "Installation instructions:"
    echo "  macOS: brew install wrk"
    echo "  Ubuntu: sudo apt-get install wrk"
    exit 1
fi

# Check if service is running
if ! curl -s http://localhost:8080/health > /dev/null
then
    echo "Service is not running on http://localhost:8080"
    echo "Please start the service before running benchmarks:"
    echo "  go run cmd/translator/main.go"
    exit 1
fi

echo "Running benchmarks..."

# Benchmark the home page
echo
echo "1. Home page (/)"
wrk -t4 -c100 -d30s http://localhost:8080/

# Benchmark the API with a simple translation request
echo
echo "2. API translation (/api/translate)"
curl -s -X POST http://localhost:8080/api/translate \
  -H "Content-Type: application/json" \
  -d '{"text":"Hello, world!","model":"gpt-3.5"}' > /tmp/translation_response.json

echo "Sample translation response:"
cat /tmp/translation_response.json | jq .
rm /tmp/translation_response.json

echo
echo "Benchmark completed."
