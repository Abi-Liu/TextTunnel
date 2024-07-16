package server

import (
	"net/http"

	"github.com/Abi-Liu/TextTunnel/config"
	"github.com/Abi-Liu/TextTunnel/internal/server/handlers"
	"github.com/Abi-Liu/TextTunnel/internal/server/middlewares"
)

func NewRouter(c *config.Config) *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("GET /health", handlers.GetHealthCheck)
	r.HandleFunc("POST /users", middlewares.ConfigMiddleware(handlers.CreateUser, c))
	r.HandleFunc("POST /login", middlewares.ConfigMiddleware(handlers.Login, c))
	r.HandleFunc("GET /ws/{roomId}", handlers.ConnectToRoom)

	return r
}
