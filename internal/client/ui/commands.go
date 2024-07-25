package ui

import (
	"context"
	"log"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type wsConnectionMsg struct {
	conn   *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

type connectionErrorMsg struct {
	err error
}

func connectToRoom(id uuid.UUID) tea.Cmd {
	return func() tea.Msg {
		header := http.Header{}
		header.Set("Authorization", "Bearer "+httpClient.AuthToken)
		dialOpts := websocket.DialOptions{HTTPHeader: header}

		ctx, cancel := context.WithCancel(context.Background())
		c, _, err := websocket.Dial(ctx, "ws://localhost:8080/ws/"+id.String(), &dialOpts)
		if err != nil {
			cancel()
			return connectionErrorMsg{err: err}
		}

		return wsConnectionMsg{conn: c, ctx: ctx, cancel: cancel}
	}
}

type receiveMsg struct {
	msg message
}

func read(conn *websocket.Conn, ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		msg := message{}
		err := wsjson.Read(ctx, conn, &msg)
		if err != nil {
			return connectionErrorMsg{err: err}
		}
		return receiveMsg{msg: msg}
	}
}

func write(conn *websocket.Conn, ctx context.Context, userId, roomId uuid.UUID, msg string) tea.Cmd {
	return func() tea.Msg {
		type params struct {
			SenderId uuid.UUID `json:"sender_id"`
			RoomId   uuid.UUID `json:"room_id"`
			Content  string    `json:"content"`
		}
		if conn == nil {
			log.Panic("CONNECTION IS NULL")
		}
		if ctx == nil {
			log.Print("CONTEXT IS NILL")
		}
		err := wsjson.Write(ctx, conn, params{
			SenderId: userId,
			RoomId:   roomId,
			Content:  msg,
		})
		if err != nil {
			return connectionErrorMsg{err: err}
		}

		return nil
	}
}
