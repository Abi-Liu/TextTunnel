package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Abi-Liu/TextTunnel/internal/client/http"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

type message struct {
	Id         uuid.UUID `json:"id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	SenderId   uuid.UUID `json:"sender_id"`
	SenderName string    `json:"sender_name"`
	RoomId     uuid.UUID `json:"room_id"`
}

type roomModel struct {
	id          uuid.UUID
	name        string
	conn        *websocket.Conn
	receiveChan chan message
	ctx         context.Context
	cancel      context.CancelFunc
	user        http.User
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	dateStyle   lipgloss.Style
	err         error
}

func newRoomModel(width, height int) roomModel {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.CharLimit = 280
	ta.SetWidth(width)
	ta.SetHeight(3)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(width, height-12)
	vp.SetContent("Welcome to the chat room!\nType a message and hit enter to send")

	return roomModel{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		dateStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		receiveChan: make(chan message, 1024),
		err:         nil,
	}
}

func (m roomModel) Init() tea.Cmd {
	// here we should establish ws connection and other room specific info
	return tea.Batch(textarea.Blink, connectToRoom(m.id))
}

func (m roomModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		cmd   tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEsc:
			return m, navigateToPage(roomListView)
		case tea.KeyEnter:
			// m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			cmd = write(m.conn, m.ctx, m.user.ID, m.id, m.textarea.Value())
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}
	}

	return m, tea.Batch(tiCmd, vpCmd, cmd)
}

func (m roomModel) formatMessages(msg message) string {
	date, err := time.Parse("2006-01-02 15:04:05.999999 -0700 MST", msg.CreatedAt.String())
	if err != nil {
		panic(err)
	}
	formattedDate := date.Format("02 January 2006 15:04:05")
	dateString := m.senderStyle.Render(formattedDate + " UTC")

	if msg.SenderId == m.user.ID {
		return fmt.Sprintf("%s %s: %s", m.dateStyle.Render(dateString), m.senderStyle.Render("You"), msg.Content)
	}
	return fmt.Sprintf("%s %s: %s", m.dateStyle.Render(dateString), m.senderStyle.Render(msg.SenderName), msg.Content)
}

func (m roomModel) View() string {
	content := fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
	helpView := helpStyle.Render("\n esc: back • enter: send • ctrl+c: quit\n")
	formatted := lipgloss.JoinVertical(lipgloss.Left, content, helpView)
	return formatted
}
