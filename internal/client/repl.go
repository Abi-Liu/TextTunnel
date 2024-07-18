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

type Message struct {
	Id        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	SenderId  uuid.UUID `json:"sender_id"`
	RoomId    uuid.UUID `json:"room_id"`
}

type Session struct {
	c       *websocket.Conn
	errChan chan error
}

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://localhost:8080/ws/"+roomId, dialOpts)
	if err != nil {
		log.Printf("Could not connect to socket: %s", err)
		return
	}

	defer c.Close(websocket.StatusGoingAway, "going away")

	session := &Session{c: c, errChan: make(chan error)}

	go session.read(ctx, c)
	for {
		s.Scan()
		msg := s.Text()
		err := session.write(ctx, c, msg)
		if err != nil {
			log.Printf("Error writing to socket: %s", err)
			return
		}
		err = <-session.errChan
		log.Printf("Error: %s", err)
		return
	}
}

func (s *Session) read(ctx context.Context, c *websocket.Conn) {
	for {
		msg := &Message{}
		err := wsjson.Read(ctx, c, msg)
		if err != nil {
			s.errChan <- err
			return
		}
		log.Print(msg.Content)
	}
}

func (s *Session) write(ctx context.Context, c *websocket.Conn, message string) error {
	err := wsjson.Write(ctx, c, message)
	return err
}
