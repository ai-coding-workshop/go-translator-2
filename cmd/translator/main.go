package main

import (
	"fmt"
	"log"
	"net/http"

	"translator-service/internal/config"
	"translator-service/internal/handlers"
	"translator-service/internal/services"
)

func main() {
	fmt.Println("Translation Service Starting...")

	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create translator service with configuration
	translatorService := services.NewTranslatorService(cfg)

	// Create handlers with dependencies
	homeHandler := handlers.NewHomeHandler()
	translateHandler := handlers.NewTranslateHandler(translatorService)
	apiHandler := handlers.NewAPIHandler(translatorService)

	// Create a new serve mux for routing
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/translate", translateHandler)
	mux.HandleFunc("/api/translate", apiHandler)

	// Serve static files
	fs := http.FileServer(http.Dir("./web/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, mux))
}
