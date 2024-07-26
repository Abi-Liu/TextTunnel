package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Abi-Liu/TextTunnel/internal/database"
	"github.com/Abi-Liu/TextTunnel/internal/server/models"
	"github.com/google/uuid"
)

func (c *Config) CreateRoom(w http.ResponseWriter, r *http.Request, user database.User) {
	type Req struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	req := &Req{}
	err := decoder.Decode(req)

	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Failed to decode parameters: %s", err))
		return
	}

	room, err := c.DB.CreateRoom(r.Context(), database.CreateRoomParams{
		ID:        uuid.New(),
		Name:      req.Name,
		CreatorID: user.ID,
		OwnerID:   user.ID,
	})

	if err != nil {
		RespondWithError(w, 400, "Failed to create room: "+err.Error())
		return
	}

	c.Hub.CreateRoom(room)

	RespondWithJson(w, 201, models.DatabaseRoomToRoom(room))
}

func (c *Config) GetAllRooms(w http.ResponseWriter, r *http.Request, _ database.User) {
	rooms, err := c.DB.FindAllRooms(r.Context())
	if err != nil {
		RespondWithError(w, 500, "Error occured when retrieving rooms: "+err.Error())
		return
	}

	RespondWithJson(w, 200, models.DatabaseRoomsToRooms(rooms))
}
