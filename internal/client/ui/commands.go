package ui

import (
	"context"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

type wsConnectionMsg struct {
	conn *websocket.Conn
}

type connectionErrorMsg struct {
	err error
}

func connectToRoom(id uuid.UUID) tea.Cmd {
	return func() tea.Msg {
		header := http.Header{}
		header.Set("Authorization", "Bearer "+httpClient.AuthToken)
		dialOpts := websocket.DialOptions{HTTPHeader: header}

		c, _, err := websocket.Dial(context.Background(), "ws://localhost:8080/ws/"+id.String(), &dialOpts)
		if err != nil {
			return connectionErrorMsg{err: err}
		}

		return wsConnectionMsg{conn: c}
	}
}
