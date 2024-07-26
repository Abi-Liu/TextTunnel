package ws

import (
	"log"
	"time"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

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
	Broadcast chan *Message
}

func (r *Room) RunRoom() {
	for {
		select {
		case client := <-r.Join:
			r.Clients[client.ID.String()] = client
			log.Printf("Client Joined: %v", client.Username)
		case client := <-r.Leave:
			delete(r.Clients, client.ID.String())
			err := client.Conn.Close(websocket.StatusNormalClosure, "Normal Closure")
			if err != nil {
				return
			}
			log.Printf("Client left: %v", client.Username)
		case msg := <-r.Broadcast:
			for _, client := range r.Clients {
				client.Receive <- msg
			}
		}
	}
}
