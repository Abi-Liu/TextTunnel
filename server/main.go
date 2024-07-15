package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	env, err := loadEnv()
	if err != nil {
		log.Fatalf("Error loading env: %v", err)
	}
	r := http.NewServeMux()

	r.HandleFunc("GET /health", getHealthCheck)

	server := &http.Server{
		Addr:              ":" + env.port,
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on port %s", env.port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getHealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, 200, "ok")
}
