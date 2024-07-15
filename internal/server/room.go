package server

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	Id        uuid.UUID
	Clients   map[*Client]struct{}
	Name      string
	CreatorId uuid.UUID
	OwnerId   uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
