package main

import (
	"log"

	"github.com/Abi-Liu/TextTunnel/internal/client/auth"
)

// import (
// 	"context"
// 	"log"
// 	"time"
//
// 	"nhooyr.io/websocket"
// 	"nhooyr.io/websocket/wsjson"
// )

func main() {
	// Dials a server, writes a single JSON message and then
	// closes the connection.

	// ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	// defer cancel()

	// c, _, err := websocket.Dial(ctx, "ws://localhost:8080/ws", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer c.CloseNow()

	// err = wsjson.Write(ctx, c, "hi")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// c.Close(websocket.StatusNormalClosure, "")

	fileSystem := auth.OSFileSystem{}
	configManager := auth.ConfigManager{FS: &fileSystem}

	_, err := configManager.LoadToken()
	if err != nil {
		log.Print("Expected err")
	}

	err = configManager.SaveToken("test")
	if err != nil {
		log.Print("Expected nil error when asving")
	}

	token, err := configManager.LoadToken()
	log.Print(token)

}
