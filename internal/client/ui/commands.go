package ui

import (
	"context"
	"log"

	"github.com/Abi-Liu/TextTunnel/internal/client/http"
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
		c, ctx, cancel, err := httpClient.ConnectToSocket(id)
		if err != nil {
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

type navigateToPageMsg struct {
	state sessionState
}

func navigateToPage(state sessionState) tea.Cmd {
	return func() tea.Msg {
		return navigateToPageMsg{state: state}
	}
}

type authorizationMsg struct {
	user http.User
}

func authorizationCmd(user http.User) tea.Cmd {
	return func() tea.Msg {
		return authorizationMsg{user: user}
	}
}

func navigateToRoom(id uuid.UUID, name string) tea.Cmd {
	return func() tea.Msg {
		return navigateToRoomMsg{
			id:   id,
			name: name,
		}
	}
}

// TODO: handle these messages inside the main model update function
type roomCreationErrorMsg struct {
	err error
}

type roomCreatedMsg struct {
	room room
}

func (m roomListModel) createRoom(name string) tea.Cmd {

	return func() tea.Msg {
		room, err := httpClient.CreateRoom(name)
		if err != nil {
			return roomCreationErrorMsg{err: err}
		}
		return roomCreatedMsg{room: httpRoomToRoom(room)}
	}
}
