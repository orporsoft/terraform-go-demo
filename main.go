package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(HealthResponse{
		Status: "ok",
	})
}

func main() {
	http.HandleFunc("/health", healthHandler)

	log.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %w", err)
	}
}
