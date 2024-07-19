package ui

import (
	"github.com/Abi-Liu/TextTunnel/internal/client/auth"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState int

const (
	UnauthorizedView = iota
	ChatListView
	RoomView
)

type MainModel struct {
	State     sessionState
	AuthToken string
	Width     int
	Height    int
}

func NewMainModel() MainModel {
	fs := auth.OSFileSystem{}
	cm := auth.ConfigManager{FS: &fs}
	token, err := cm.LoadToken()
	var state sessionState
	if err != nil {
		token = ""
		state = 0
	}

	return MainModel{
		State:     state,
		AuthToken: token,
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.Height = msg.Height
		m.Width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil

}

func (m MainModel) View() string {
	return m.AuthToken
}
