package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func navigateToPage(state sessionState) tea.Cmd {
	return func() tea.Msg {
		return navigateToPageMsg{state: state}
	}
}

type UnauthorizedModel struct {
	focusIndex int
}

func NewUnauthorizedModel() UnauthorizedModel {
	return UnauthorizedModel{}
}

func (m UnauthorizedModel) Init() tea.Cmd {
	return nil
}

func (m UnauthorizedModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			if m.focusIndex > 0 {
				m.focusIndex--
			}
		case "l", "right":
			if m.focusIndex < 1 {
				m.focusIndex++
			}
		case "enter":
			// create a new command to route to the appropriate page
			if m.focusIndex == 0 {
				cmd = navigateToPage(loginView)
			} else {
				cmd = navigateToPage(signUpView)
			}
		}
	}
	return m, cmd
}

func (m UnauthorizedModel) View() string {
	focusedLoginButton := focusedStyle.Render("[ Login ]")
	blurredLoginButton := fmt.Sprintf("[ %s ]", blurredStyle.Render("Login"))
	focusedSignupButton := focusedStyle.Render("[ Sign up ]")
	blurredSignupButton := fmt.Sprintf("[ %s ]", blurredStyle.Render("Sign up"))

	login := &blurredLoginButton
	if m.focusIndex == 0 {
		login = &focusedLoginButton
	}

	signup := &blurredSignupButton
	if m.focusIndex == 1 {
		signup = &focusedSignupButton
	}
	s := "Welcome to TextTunnel\n\n"
	s += fmt.Sprintf("   %s   %s   ", *login, *signup)
	return s
}
