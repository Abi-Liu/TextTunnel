package server

import (
	"log"
	"net/http"

	"github.com/Abi-Liu/TextTunnel/internal/database"
	"github.com/Abi-Liu/TextTunnel/internal/server/ws"
	"nhooyr.io/websocket"
)

func (c *Config) ConnectToRoom(w http.ResponseWriter, r *http.Request, user database.User) {
	roomId := r.PathValue("roomId")
	room, ok := c.Hub.Rooms[roomId]
	if !ok {
		RespondWithError(w, 404, "Room does not exist")
		log.Printf("Error, room does not exist")
		return
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Printf("error upgrading connection: %s", err)
		return
	}

	defer conn.CloseNow()

	client := &ws.Client{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Receive:   make(chan *ws.Message, 1024),
		Room:      room,
		Conn:      conn,
		DB:        c.DB,
	}

	room.Join <- client

	go client.Read()
	client.Write()

	conn.Close(websocket.StatusNormalClosure, "")
}
