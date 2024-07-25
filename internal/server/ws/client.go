package ws

import (
	"context"
	"log"
	"strings"
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
	Ctx       context.Context
}

type Message struct {
	Id        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	SenderId  uuid.UUID `json:"sender_id"`
	Username  string    `json:"username"`
	RoomId    uuid.UUID `json:"room_id"`
}

func databaseMessageToMessage(message database.Message, username string) Message {
	return Message{
		Id:        message.ID,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
		SenderId:  message.SenderID,
		Username:  username,
		RoomId:    message.RoomID,
	}
}

func (c *Client) Write() {
	for {
		select {
		case msg := <-c.Receive:
			err := wsjson.Write(c.Ctx, c.Conn, msg)
			if err != nil {
				if !strings.Contains(err.Error(), "closed network connection") {
					log.Printf("Failed to write to client: %s", err)
				}
				return
			}

		case <-c.Ctx.Done():
			log.Printf("Context done in write method")
			return
		}
	}
}

func (c *Client) Read() {
	for {
		select {
		case <-c.Ctx.Done():
			log.Printf("Context done")
			return
		default:
			type Req struct {
				SenderId uuid.UUID `json:"sender_id"`
				RoomId   uuid.UUID `json:"room_id"`
				Content  string    `json:"content"`
			}
			req := &Req{}
			err := wsjson.Read(c.Ctx, c.Conn, req)
			if err != nil {
				if websocket.CloseStatus(err) != websocket.StatusNormalClosure {
					log.Printf("Error reading from client: %s\n%v", err, req)
				}
				return
			}

			dbMsg, err := c.DB.CreateMessage(c.Ctx, database.CreateMessageParams{
				ID:       uuid.New(),
				Content:  req.Content,
				SenderID: c.ID,
				RoomID:   c.Room.ID,
			})
			if err != nil {
				log.Printf("Failed to insert message: %s", err)
				return
			}

			msg := databaseMessageToMessage(dbMsg, c.Username)
			c.Room.Broadcast <- &msg
		}
	}
}
