package ws

import (
	"context"
	"log"
	"time"

	"github.com/Abi-Liu/TextTunnel/internal/database"
	"github.com/google/uuid"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Client struct {
	ID        uuid.UUID
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Conn      *websocket.Conn
	Receive   chan *Message
	Room      *Room
	DB        *database.Queries
}

type Message struct {
	Id        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	SenderId  uuid.UUID `json:"sender_id"`
	RoomId    uuid.UUID `json:"room_id"`
}

func databaseMessageToMessage(message database.Message) Message {
	return Message{
		Id:        message.ID,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
		SenderId:  message.SenderID,
		RoomId:    message.RoomID,
	}
}

func (c *Client) Write() {
	for msg := range c.Receive {
		wsjson.Write(context.Background(), c.Conn, msg)
	}
}

func (c *Client) Read() {
	type Req struct {
		Content  string    `json:"Content"`
		SenderId uuid.UUID `json:"sender_id"`
		RoomId   uuid.UUID `json:"room_id"`
	}
	for {
		req := &Req{}
		err := wsjson.Read(context.Background(), c.Conn, req)
		if err != nil {
			log.Printf("Error reading from client: %s", err)
		}
		dbMsg, err := c.DB.CreateMessage(context.Background(), database.CreateMessageParams{
			ID:       c.ID,
			Content:  req.Content,
			SenderID: req.SenderId,
			RoomID:   req.RoomId,
		})
		if err != nil {
			log.Printf("Failed to insert message: %s", err)
		}

		msg := databaseMessageToMessage(dbMsg)
		c.Room.Broadcast <- &msg
	}
}
