package client

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"os"

	"github.com/Abi-Liu/TextTunnel/internal/client/auth"
	"github.com/google/uuid"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const ROOM_ID = "ea4d1f5a-8e07-4588-b7be-43c563960590"

func Start(cm *auth.ConfigManager) {
	fmt.Println("Welcome to TextTunnel")

	// ONLY USING FOR TESTING PURPOSES
	token := "75a23f0f41f47f739eed7baf384811b3749b177da377752717d067d74fdac45c"
	roomId := "ea4d1f5a-8e07-4588-b7be-43c563960590"
	client := CreateHttpClient()
	client.SetAuthToken(token)
	s := bufio.NewScanner(os.Stdin)

	header := http.Header{}
	header.Set("Authorization", "Bearer "+token)
	dialOpts := &websocket.DialOptions{HTTPHeader: header}
	c, _, err := websocket.Dial(context.Background(), "ws://localhost:8080/ws/"+roomId, dialOpts)
	if err != nil {
		log.Printf("Could not connect to socket: %s", err)
		return
	}

	go read(c)
	for {
		s.Scan()
		msg := s.Text()
		write(c, msg)
	}
}

func read(c *websocket.Conn) {
	type Message struct {
		Id        uuid.UUID `json:"id"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		SenderId  uuid.UUID `json:"sender_id"`
		RoomId    uuid.UUID `json:"room_id"`
	}

	for {
		msg := &Message{}
		wsjson.Read(context.Background(), c, msg)
		log.Print(msg)

	}
}

func write(c *websocket.Conn, message string) {
	wsjson.Write(context.Background(), c, message)
}
