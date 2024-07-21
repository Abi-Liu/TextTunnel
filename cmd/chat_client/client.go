package main

import (
	"log"

	"github.com/Abi-Liu/TextTunnel/internal/client/auth"
	"github.com/Abi-Liu/TextTunnel/internal/client/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fs := auth.OSFileSystem{}
	cm := auth.ConfigManager{FS: &fs}
	token, err := cm.LoadToken()
	if err != nil {
		token = ""
	}
	model := ui.NewMainModel(token)
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}
