package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Abi-Liu/TextTunnel/internal/database"
	"github.com/Abi-Liu/TextTunnel/internal/server/auth"
	"github.com/Abi-Liu/TextTunnel/internal/server/models"
	"github.com/google/uuid"
)

func (c *Config) Login(w http.ResponseWriter, r *http.Request) {
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

	user, err := c.DB.FindUserByUsername(r.Context(), req.Username)
	if err != nil {
		RespondWithError(w, 404, fmt.Sprintf("User %s not found", req.Username))
		return
	}

	ok := auth.CompareHashAndPassword(user.Password, req.Password)
	if !ok {
		RespondWithError(w, 400, "Unauthorized - Passwords do not match")
		return
	}

	RespondWithJson(w, 200, models.DatabaseUserToUser(user))
}

func (c *Config) CreateUser(w http.ResponseWriter, r *http.Request) {
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
