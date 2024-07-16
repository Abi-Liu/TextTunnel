package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Abi-Liu/TextTunnel/config"
	"github.com/Abi-Liu/TextTunnel/internal/server"
)

func main() {
	config, err := config.CreateConfig()
	if err != nil {
		log.Fatalf("Error creating config: %v", err)
	}

	r := server.NewRouter(config)

	server := &http.Server{
		Addr:              ":" + config.Env.PORT,
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on port %s", config.Env.PORT)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
