package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type inputModel struct {
	input textinput.Model
}

func newInputModel() tea.Model {
	m := inputModel{input: textinput.New()}
	m.input.Cursor.Style = cursorStyle
	m.input.CharLimit = 24
	m.input.Placeholder = "Enter room name"
	m.input.Focus()
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
	return m.input.View()
}
