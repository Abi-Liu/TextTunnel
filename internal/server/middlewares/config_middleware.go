package middlewares

import (
	"net/http"

	"github.com/Abi-Liu/TextTunnel/config"
)

type ConfigHandler func(http.ResponseWriter, *http.Request, *config.Config)

func ConfigMiddleware(handler ConfigHandler, c *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, c)
	})
}
