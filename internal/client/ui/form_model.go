package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

var focusedButton string
var blurredButton string

type FormModel struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
	error      string
	state      sessionState
}

func NewFormModel(state sessionState) FormModel {
	m := FormModel{state: state}
	if state == loginView {
		m.inputs = make([]textinput.Model, 2)
		focusedButton = focusedStyle.Render("[ Login ]")
		blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Login"))
	} else {
		m.inputs = make([]textinput.Model, 3)
	}

	for i := range m.inputs {
		input := textinput.New()
		input.Cursor.Style = cursorStyle
		input.CharLimit = 24

		switch i {
		case 0:
			input.Placeholder = "Username"
			input.Focus()
			input.PromptStyle = focusedStyle
			input.TextStyle = focusedStyle
		case 1:
			input.Placeholder = "Password"
			input.EchoMode = textinput.EchoPassword
			input.EchoCharacter = '*'
		case 2:
			input.Placeholder = "Confirm Password"
			input.EchoMode = textinput.EchoPassword
			input.EchoCharacter = '*'

		}
		m.inputs[i] = input
	}

	return m
}

func (m FormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m FormModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()
			// exit if the user presses enter while the submit button was focused
			if s == "enter" && m.focusIndex == len(m.inputs) {
				valid := validateInputs(m.inputs)
				if !valid {
					m.error = "Username and password cannot be empty\n"
				}
				if m.state == signUpView {
					valid = validatePassword(m.inputs)
					if !valid {
						m.error = "Passwords do not match"
					}
				}
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(message)
	return m, cmd
}

func (m *FormModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m FormModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	if m.error != "" {
		fmt.Fprintf(&b, "%s\n\n", m.error)
	}

	return b.String()
}

func validateInputs(inputs []textinput.Model) bool {
	for _, v := range inputs {
		if len(v.Value()) == 0 {
			return false
		}
	}

	return true
}

func validatePassword(inputs []textinput.Model) bool {
	return inputs[1].Value() == inputs[2].Value()
}
