package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type HealthResponse struct {
	Status string `json:"status"`
	Version string `json:"version"`
	Env string `json:"env"`
	APIKEYConfigured bool `json:"api_key_configured"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	env := os.Getenv("APP_ENV")
	apiKey := os.Getenv("API_KEY")

	if env == "" {
		env = "unknown"
	}

	response := HealthResponse{
		Status: "ok",
		Version: "1.2",
		Env: env,
		APIKEYConfigured: apiKey != "",
	}

	w.Header().Set("Content-Type", "application/json")

	if err :=json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func main() {
	http.HandleFunc("/health", healthHandler)

	log.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
