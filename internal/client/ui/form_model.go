package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/Abi-Liu/TextTunnel/internal/client/http"
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
		focusedButton = focusedStyle.Render("[ Sign up ]")
		blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Sign up"))
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
					user, err := signUp(m.inputs[0].Value(), m.inputs[1].Value())
					if err != nil {
						m.error = err.Error()
						return m, nil
					}
					cmd = authorizationCmd(user)
				} else {
					user, err := login(m.inputs[0].Value(), m.inputs[1].Value())
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

func signUp(username, password string) (http.User, error) {
	type Req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	req := Req{
		Username: username,
		Password: password,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return http.User{}, err
	}

	reqBody := bytes.NewReader(data)
	res, err := httpClient.Post(http.BASE_URL+"/users", reqBody)
	if err != nil {
		return http.User{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	var user http.User
	var error http.Error

	if res.StatusCode >= 400 {
		err = json.Unmarshal(body, &error)
		if err != nil {
			return http.User{}, err
		}
		return http.User{}, fmt.Errorf(error.Error)
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return http.User{}, err
	}
	return user, nil
}

func login(username, password string) (http.User, error) {
	type Req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	fmt.Println(username, password)

	req := Req{
		Username: username,
		Password: password,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return http.User{}, err
	}

	reqBody := bytes.NewReader(data)
	res, err := httpClient.Post(http.BASE_URL+"/login", reqBody)
	if err != nil {
		return http.User{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	var user http.User
	var error http.Error
	if res.StatusCode >= 400 {
		err = json.Unmarshal(body, &error)
		if err != nil {
			return http.User{}, err
		}
		return http.User{}, fmt.Errorf(error.Error)
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return http.User{}, err
	}
	return user, nil
}

func authorizationCmd(user http.User) tea.Cmd {
	return func() tea.Msg {
		return authorizationMsg{user: user}
	}
}
