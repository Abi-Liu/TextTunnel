package server

import (
	"net/http"

	"github.com/Abi-Liu/TextTunnel/internal/database"
	"github.com/Abi-Liu/TextTunnel/internal/server/auth"
)

type AuthHandler func(http.ResponseWriter, *http.Request, database.User)

func (c Config) EnsureAuth(handler AuthHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key, err := auth.GetAuthorizationKey(r)
		if err != nil {
			RespondWithError(w, 401, "Invalid authorization key")
			return
		}

		user, err := c.DB.FindUserByApiKey(r.Context(), key)
		if err != nil {
			RespondWithError(w, 401, "Invalid authorization key")
			return
		}

		handler(w, r, user)
	})
}
