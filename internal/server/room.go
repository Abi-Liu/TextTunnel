package server

import (
	"github.com/Abi-Liu/TextTunnel/internal/database"
	"github.com/google/uuid"
)

type Room struct {
	Id        uuid.UUID
	Clients   map[*Client]struct{}
	Name      string
	CreatorId uuid.UUID
	OwnerId   uuid.UUID
	Forward   chan *Message
	Join      chan *Client
	Leave     chan *Client
}

func NewRoom(dbRoom database.Room) *Room {
	return &Room{
		Id:        dbRoom.ID,
		Clients:   make(map[*Client]struct{}),
		Name:      dbRoom.Name,
		CreatorId: dbRoom.CreatorID,
		OwnerId:   dbRoom.OwnerID,
		Forward:   make(chan *Message),
		Join:      make(chan *Client),
		Leave:     make(chan *Client),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = struct{}{}
		case client := <-r.Leave:
			delete(r.Clients, client)
			close(client.receive)
		case message := <-r.Forward:
			for client := range r.Clients {
				client.receive <- message
			}
		}
	}
}
