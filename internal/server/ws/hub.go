package ws

import (
	"context"

	"github.com/Abi-Liu/TextTunnel/internal/database"
)

type Hub struct {
	Rooms map[string]*Room
}

func CreateHub(db *database.Queries) *Hub {
	rooms, err := db.FindAllRooms(context.Background())

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

// for each room we run the room on a new goroutine
func (h *Hub) Run() {
	for _, room := range h.Rooms {
		go room.RunRoom()
	}
}

func (h *Hub) CreateRoom(r database.Room) bool {
	_, ok := h.Rooms[r.ID.String()]
	if ok {
		return false
	}

	room := &Room{
		ID:        r.ID,
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		OwnerId:   r.OwnerID,
		CreatorId: r.CreatorID,
		Clients:   make(map[string]*Client),
		Join:      make(chan *Client, 1024),
		Leave:     make(chan *Client, 1024),
		Broadcast: make(chan *Message, 1024),
	}
	h.Rooms[r.ID.String()] = room

	go room.RunRoom()
	return true
}
