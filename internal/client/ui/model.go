package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState int

const (
	unauthorizedView sessionState = iota
	loginView
	signUpView
	roomListView
	roomView
)

type MainModel struct {
	State             sessionState
	AuthToken         string
	Width             int
	Height            int
	UnauthorizedModel tea.Model
	LoginModel        tea.Model
}

type navigateToPageMsg struct {
	state sessionState
}

func NewMainModel(token string) MainModel {
	state := roomListView
	if token == "" {
		state = unauthorizedView
	}

	login := NewLoginModel()
	unauthorized := NewUnauthorizedModel()

	return MainModel{
		State:             state,
		AuthToken:         token,
		LoginModel:        login,
		UnauthorizedModel: unauthorized,
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.Height = msg.Height
		m.Width = msg.Width
	case navigateToPageMsg:
		m.State = msg.state
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			switch m.State {
			case unauthorizedView:
				unauthorized, command := m.UnauthorizedModel.Update(msg)
				m.UnauthorizedModel = unauthorized
				cmd = command
			case loginView:
				login, command := m.LoginModel.Update(message)
				m.LoginModel = login
				cmd = command
			}
		}
	}
	return m, cmd

}

func (m MainModel) View() string {
	switch m.State {
	case unauthorizedView:
		// show unauthorized view
		return m.UnauthorizedModel.View()
	case loginView:
		// show login view
		return m.LoginModel.View()
	case signUpView:
		// show sign up view
	case roomListView:
		// show list of rooms available
	case roomView:
		// show the room view with chat etc.
	}
	return ""
}
