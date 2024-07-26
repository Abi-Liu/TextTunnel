package ui

import (
	"time"

	"github.com/Abi-Liu/TextTunnel/internal/client/http"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type navigateToRoomMsg struct {
	id   uuid.UUID
	name string
}

type room struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	OwnerId   uuid.UUID
	CreatorId uuid.UUID
}

func (r room) FilterValue() string {
	return r.Name
}

func (r room) Title() string {
	return r.Name
}

func (r room) Description() string {
	return ""
}

type roomListModel struct {
	list       list.Model
	input      tea.Model
	showInput  bool
	err        error
	focusIndex int
}

type populateListMsg struct {
	list []list.Item
}

func newRoomListModel() *roomListModel {
	list := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	list.SetShowHelp(false)
	input := newInputModel()
	return &roomListModel{
		list:      list,
		input:     input,
		showInput: false,
	}
}

// fetch the rooms and populate the list
func (m *roomListModel) initList(width, height int) {
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height-3)
	m.list.SetShowHelp(false)
}

func httpRoomToRoom(r http.Room) room {
	return room{
		ID:        r.ID,
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		OwnerId:   r.OwnerId,
		CreatorId: r.CreatorId,
	}
}

func (m roomListModel) Init() tea.Cmd {
	return func() tea.Msg {
		rooms, err := httpClient.FetchRooms()
		if err != nil {
			m.err = err
			return nil
		}
		list := make([]list.Item, len(rooms))
		for i, r := range rooms {
			list[i] = httpRoomToRoom(r)
		}
		return populateListMsg{list: list}
	}
}

func (m roomListModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if !m.showInput {
				// handle navigating to the correct room.
				room := m.list.Items()[m.focusIndex].(room)
				cmd = navigateToRoom(room.ID, room.Name)
			} else {
				// create a new room command
			}
		case "up", "k":
			if m.focusIndex > 0 {
				m.focusIndex--
				m.list.Select(m.focusIndex)
			}
		case "down", "j":
			if m.focusIndex < len(m.list.Items())-1 {
				m.focusIndex++
				m.list.Select(m.focusIndex)
			}
		case "c":
			m.showInput = true
			cmd = m.input.Init()
		case "esc":
			m.showInput = false
		default:
			if m.showInput {
				m.input, cmd = m.input.Update(msg)
			}
		}
	}
	return m, cmd
}

func (m roomListModel) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	helpView := helpStyle.Render("\n c: create room • enter: join room • /: filter • ctrl+c: quit\n")
	s := lipgloss.JoinVertical(lipgloss.Left, m.list.View(), helpView)
	if m.showInput {
		s = lipgloss.JoinVertical(lipgloss.Left, s, m.input.View())
	}
	return s
}
