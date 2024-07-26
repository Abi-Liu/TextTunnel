package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type inputModel struct {
	input textinput.Model
	err   error
}

func newInputModel() tea.Model {
	m := inputModel{input: textinput.New()}
	m.input.Cursor.Style = cursorStyle
	m.input.CharLimit = 24
	m.input.Placeholder = "Enter room name"
	m.input.PromptStyle = focusedStyle
	m.input.TextStyle = focusedStyle
	return m
}

func (m inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m inputModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := message.(type) {
	case tea.KeyMsg:
		m.input, cmd = m.input.Update(msg)

	}
	return m, cmd
}

func (m inputModel) View() string {
	if m.err != nil {
		return lipgloss.JoinVertical(lipgloss.Left, "Could not create room", m.input.View())
	}
	return m.input.View()
}
