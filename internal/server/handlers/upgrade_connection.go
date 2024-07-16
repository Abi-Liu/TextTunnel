package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func ConnectToRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Printf("error upgrading connection: %s", err)
		return
	}

	defer conn.CloseNow()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var v interface{}
	err = wsjson.Read(ctx, conn, &v)

	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("%v", v)

	conn.Close(websocket.StatusNormalClosure, "")
}
