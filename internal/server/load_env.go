package server

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	PORT   string
	DB_URL string
}

func LoadEnv() (*Env, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("port env missing")
	}

	db := os.Getenv("DATABASE_URL")
	if db == "" {
		return nil, errors.New("db string missing")
	}

	return &Env{
		PORT:   port,
		DB_URL: db,
	}, nil
}
