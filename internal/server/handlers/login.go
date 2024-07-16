package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Abi-Liu/TextTunnel/config"
	"github.com/Abi-Liu/TextTunnel/internal/database"
	"github.com/Abi-Liu/TextTunnel/internal/server/auth"
	"github.com/Abi-Liu/TextTunnel/internal/server/models"
	"github.com/google/uuid"
)

func Login(w http.ResponseWriter, r *http.Request, c *config.Config) {

}

func CreateUser(w http.ResponseWriter, r *http.Request, c *config.Config) {
	type Req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	req := &Req{}
	err := decoder.Decode(req)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Could not decode parameters: %s", err))
		return
	}

	hashedPw, err := auth.HashPassword(req.Password)
	if err != nil {
		RespondWithError(w, 500, "Failed to hash password")
		return
	}

	user, err := c.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:       uuid.New(),
		Username: req.Username,
		Password: hashedPw,
	})

	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			RespondWithError(w, 400, "Username already exists, Please choose another")
			return
		}
		RespondWithError(w, 500, "Failed to create user")
		return
	}

	RespondWithJson(w, http.StatusCreated, models.DatabaseUserToUser(user))
}
