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
	RoomModel         tea.Model
}

type navigateToPageMsg struct {
	state sessionState
}

type authorizationMsg struct {
	user http.User
}

type validAuthTokenOnStartMsg struct {
	user http.User
}

type navigateToRoomsListMsg struct{}

var httpClient *http.HttpClient

func NewMainModel(token string, cm auth.ConfigManager) MainModel {
	httpClient = http.CreateHttpClient(token)

	unauthorized := NewUnauthorizedModel()
	login := NewFormModel(loginView)
	signUp := NewFormModel(signUpView)
	roomList := newRoomListModel()
	room := newRoomModel(30, 5)

	return MainModel{
		State:             unauthorizedView,
		Cm:                cm,
		AuthToken:         token,
		LoginModel:        login,
		UnauthorizedModel: unauthorized,
		SignUpModel:       signUp,
		RoomListModel:     roomList,
		RoomModel:         room,
	}
}

func (m MainModel) Init() tea.Cmd {
	var cmd tea.Cmd
	if m.AuthToken != "" {
		user, err := httpClient.GetUserByAuthToken()
		if err != nil {
			panic("Corrupted auth token!")
		}

		cmd = func() tea.Msg {
			return validAuthTokenOnStartMsg{
				user: user,
			}
		}
	}
	return cmd
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
		m.RoomModel = newRoomModel(msg.Width, msg.Height)
	case navigateToPageMsg:
		m.State = msg.state
	case validAuthTokenOnStartMsg:
		m.State = roomListView
		m.User = msg.user
		cmd = m.RoomListModel.Init()
	case authorizationMsg:
		m.User = msg.user
		// save the authorization token
		m.Cm.SaveToken(msg.user.ApiKey)
		httpClient.SetAuthToken(msg.user.ApiKey)

		// switch the state to the room list view
		m.State = roomListView
		// initialize the list
		cmd = m.RoomListModel.Init()
	case populateListMsg:
		model := m.RoomListModel.(*roomListModel)
		model.list.SetItems(msg.list)
		m.RoomListModel = model
	case navigateToRoomMsg:
		model := m.RoomModel.(roomModel)
		model.name = msg.name
		model.id = msg.id
		m.RoomModel = model
		m.State = roomView
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
			case roomView:
				m.RoomModel, cmd = m.RoomModel.Update(message)
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
		return m.RoomModel.View()
	}
	return ""
}
