package ws

import (
	"time"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

type Client struct {
	ID        uuid.UUID
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Conn      *websocket.Conn
	Receive   chan *Message
	Room      *Room
}

type Message struct {
	Id        uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	SenderId  uuid.UUID
	RoomId    uuid.UUID
}
