package handlers

import (
	"net/http"

	"github.com/Abi-Liu/TextTunnel/internal/server"
)

func GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	server.RespondWithJson(w, 200, "OK")
}
