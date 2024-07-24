package server

import (
	"net/http"
)

func (c *Config) NewRouter() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("GET /health", GetHealthCheck)
	r.HandleFunc("POST /users", c.CreateUser)
	r.HandleFunc("POST /login", c.Login)
	r.HandleFunc("GET /users", c.EnsureAuth(c.GetUser))

	r.HandleFunc("GET /ws/{roomId}", c.EnsureAuth(c.ConnectToRoom))
	r.HandleFunc("POST /rooms", c.EnsureAuth(c.CreateRoom))
	r.HandleFunc("GET /rooms", c.EnsureAuth(c.GetAllRooms))

	return r
}
