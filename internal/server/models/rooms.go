package models

import (
	"time"

	"github.com/Abi-Liu/TextTunnel/internal/database"
	"github.com/google/uuid"
)

type Room struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatorId uuid.UUID `json:"creator_id"`
	OwnerId   uuid.UUID `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DatabaseRoomToRoom(r database.Room) Room {
	return Room{
		ID:        r.ID,
		Name:      r.Name,
		CreatorId: r.CreatorID,
		OwnerId:   r.OwnerID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func DatabaseRoomsToRooms(rooms []database.Room) []Room {
	res := make([]Room, len(rooms))
	for i, room := range rooms {
		res[i] = DatabaseRoomToRoom(room)
	}
	return res
}
