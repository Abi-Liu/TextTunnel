package server

import (
	"net/http"

	"github.com/Abi-Liu/TextTunnel/config"
	"github.com/Abi-Liu/TextTunnel/internal/server/handlers"
	"github.com/Abi-Liu/TextTunnel/internal/server/ws"
)

func NewRouter(c *config.Config, hub *ws.Hub) *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("GET /health", handlers.GetHealthCheck)
	r.HandleFunc("GET /ws/{roomId}", handlers.ConnectToRoom)

	return r
}
