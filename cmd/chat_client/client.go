package main

import (
	"log"

	"github.com/Abi-Liu/TextTunnel/internal/client/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := ui.NewMainModel()
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}
