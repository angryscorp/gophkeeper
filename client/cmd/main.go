package main

import (
	"fmt"
	"gophkeeper/client/internal/tui/menu"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	program := tea.NewProgram(menu.New(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Something went wrong: %v", err)
		os.Exit(1)
	}
}
