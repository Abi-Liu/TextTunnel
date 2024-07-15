package server

import (
	"time"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

type Client struct {
	conn    *websocket.Conn
	receive chan *Message
	room    *Room
}

type Message struct {
	Id        uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	SenderId  uuid.UUID
	RoomId    uuid.UUID
}
