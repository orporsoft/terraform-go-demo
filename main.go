package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthResponse struct {
	Status           string `json:"status"`
	Version          string `json:"version"`
	Env              string `json:"env"`
	APIKeyConfigured bool   `json:"api_key_configured"`
	Database         string `json:"database"`
}

var db *pgxpool.Pool

func healthHandler(w http.ResponseWriter, r *http.Request) {
	env := os.Getenv("APP_ENV")
	apiKey := os.Getenv("API_KEY")

	if env == "" {
		env = "unknown"
	}

	databaseStatus := "DOWN"

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err == nil {
		databaseStatus = "UP"
	}

	response := HealthResponse{
		Status:           "ok",
		Version:          "1.6",
		Env:              env,
		APIKeyConfigured: apiKey != "",
		Database:         databaseStatus,
	}

	w.Header().Set("Content-Type", "application/json")

	if databaseStatus != "UP" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func main() {
	ctx := context.Background()

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	if dbHost == "" {
		dbHost = "localhost"
	}

	if dbPort == "" {
		dbPort = "5432"
	}

	databaseURL := "postgres://" +
		dbUser + ":" +
		dbPassword + "@" +
		dbHost + ":" +
		dbPort + "/" +
		dbName

	log.Printf("Connecting to PostgreSQL at %s:%s/%s", dbHost, dbPort, dbName)

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to create PostgreSQL pool: %v", err)
	}

	db = pool

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Ping(pingCtx); err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	log.Println("Connected to PostgreSQL")

	defer db.Close()

	http.HandleFunc("/health", healthHandler)

	log.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}