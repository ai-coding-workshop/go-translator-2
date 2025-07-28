package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Translation Service Starting...")

	// Placeholder for future HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Translation Service is running!")
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}