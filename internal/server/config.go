package server

import (
	"database/sql"

	"github.com/Abi-Liu/TextTunnel/internal/database"
	"github.com/Abi-Liu/TextTunnel/internal/server/ws"
	_ "github.com/lib/pq"
)

type Config struct {
	DB  *database.Queries
	Hub *ws.Hub
	Env *Env
}

func CreateConfig() (*Config, error) {
	env, err := LoadEnv()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", env.DB_URL)
	if err != nil {
		return nil, err
	}

	dbQueries := database.New(db)

	hub := ws.CreateHub(dbQueries)

	return &Config{
		DB:  dbQueries,
		Hub: hub,
		Env: env,
	}, nil
}
