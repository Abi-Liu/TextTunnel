package ui

import (
	"github.com/Abi-Liu/TextTunnel/internal/client/auth"
	"github.com/Abi-Liu/TextTunnel/internal/client/http"
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
	Cm                auth.ConfigManager
	AuthToken         string
	User              http.User
	Width             int
	Height            int
	UnauthorizedModel tea.Model
	LoginModel        tea.Model
	SignUpModel       tea.Model
	RoomListModel     tea.Model
}

type navigateToPageMsg struct {
	state sessionState
}

type authorizationMsg struct {
	user http.User
}

var httpClient *http.HttpClient

func NewMainModel(token string, cm auth.ConfigManager) MainModel {
	state := roomListView
	if token == "" {
		state = unauthorizedView
	}
	httpClient = http.CreateHttpClient(token)

	unauthorized := NewUnauthorizedModel()
	login := NewFormModel(loginView)
	signUp := NewFormModel(signUpView)
	roomList := newRoomListModel()

	return MainModel{
		State:             state,
		Cm:                cm,
		AuthToken:         token,
		LoginModel:        login,
		UnauthorizedModel: unauthorized,
		SignUpModel:       signUp,
		RoomListModel:     roomList,
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
		if roomList, ok := m.RoomListModel.(*roomListModel); ok {
			roomList.initList(msg.Width, msg.Height)
		}
	case navigateToPageMsg:
		m.State = msg.state
	case authorizationMsg:
		m.User = msg.user
		// save the authorization token
		m.Cm.SaveToken(msg.user.ApiKey)
		httpClient.SetAuthToken(msg.user.ApiKey)

		// switch the state to the room list view
		m.State = roomListView
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
			case signUpView:
				signup, command := m.SignUpModel.Update(message)
				m.SignUpModel = signup
				cmd = command
			case roomListView:
				roomList, command := m.RoomListModel.Update(message)
				m.RoomListModel = roomList
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
		return m.SignUpModel.View()
	case roomListView:
		// show list of rooms available
		return m.RoomListModel.View()
	case roomView:
		// show the room view with chat etc.
	}
	return ""
}
