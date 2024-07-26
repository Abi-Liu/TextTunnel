package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()
)

type FormModel struct {
	focusIndex    int
	inputs        []textinput.Model
	error         string
	state         sessionState
	focusedButton string
	blurredButton string
}

func NewFormModel(state sessionState) FormModel {
	m := FormModel{state: state}
	if state == loginView {
		m.inputs = make([]textinput.Model, 2)
		m.focusedButton = focusedStyle.Render("[ Login ]")
		m.blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Login"))
	} else {
		m.inputs = make([]textinput.Model, 3)
		m.focusedButton = focusedStyle.Render("[ Sign up ]")
		m.blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Sign up"))
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
	var cmd tea.Cmd
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

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
				if !valid {
					return m, nil
				}
				var cmd tea.Cmd
				if m.state == signUpView {
					user, err := httpClient.PostSignUp(m.inputs[0].Value(), m.inputs[1].Value())
					if err != nil {
						m.error = err.Error()
						return m, nil
					}
					cmd = authorizationCmd(user)
				} else {
					user, err := httpClient.Login(m.inputs[0].Value(), m.inputs[1].Value())
					if err != nil {
						m.error = err.Error()
						return m, nil
					}
					cmd = authorizationCmd(user)
				}
				return m, cmd
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

		case "esc":
			cmd = navigateToPage(unauthorizedView)
			return m, cmd
		}
	}

	cmd = m.updateInputs(message)
	return m, cmd
}

func (m *FormModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further fmtic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m FormModel) View() string {
	var b strings.Builder
	helpView := helpStyle.Render("\n ↑/shift+tab: navigate up • ↓/tab: navigate down • esc: back • enter: submit • ctrl+c: quit\n")
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &m.blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &m.focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	if m.error != "" {
		fmt.Fprintf(&b, "%s\n\n", m.error)
	}

	return lipgloss.JoinVertical(lipgloss.Left, b.String(), helpView)
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
	return inputs[1].Value() == inputs[2].Value() && len(inputs[1].Value()) > 0
}
