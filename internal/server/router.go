package server

import (
	"net/http"
)

func (c *Config) NewRouter() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("GET /health", GetHealthCheck)
	r.HandleFunc("POST /users", c.CreateUser)
	r.HandleFunc("POST /login", c.Login)
	r.HandleFunc("GET /ws/{roomId}", c.ConnectToRoom)

	return r
}
