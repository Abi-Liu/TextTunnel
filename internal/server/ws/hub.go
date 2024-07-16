package ws

import (
	"context"
	"time"

	"github.com/Abi-Liu/TextTunnel/config"
	"github.com/Abi-Liu/TextTunnel/internal/database"
	"github.com/google/uuid"
)

type Hub struct {
	Rooms map[string]*Room
}

func CreateHub(c *config.Config) *Hub {
	rooms, err := c.DB.FindAllRooms(context.Background())

	hub := &Hub{
		Rooms: make(map[string]*Room),
	}
	if err != nil {
		return hub
	}
	for _, room := range rooms {
		hub.CreateRoom(room)
	}

	return hub
}

func (h *Hub) CreateRoom(r database.Room) bool {
	_, ok := h.Rooms[r.ID.String()]
	if !ok {
		return false
	}
	h.Rooms[r.ID.String()] = &Room{
		ID:        r.ID,
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		OwnerId:   r.OwnerID,
		CreatorId: r.CreatorID,
		Join:      make(chan *Client),
		Leave:     make(chan *Client),
		Forward:   make(chan *Message),
	}

	return true
}

type Room struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatorId uuid.UUID
	OwnerId   uuid.UUID
	Clients   map[string]*Client
	Join      chan *Client
	Leave     chan *Client
	Forward   chan *Message
}

type Client struct {
	id        uuid.UUID
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
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
