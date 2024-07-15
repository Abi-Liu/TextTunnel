package routes

import (
	"net/http"

	"github.com/Abi-Liu/TextTunnel/internal/server/handlers"
)

func NewRouter() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("GET /health", handlers.GetHealthCheck)

	return r
}
