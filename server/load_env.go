package main

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	port   string
	DB_URL string
}

func loadEnv() (*Env, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("Port env missing")
	}

	return &Env{
		port: port,
	}, nil
}
