package help

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	help func() string
}

func New(help func() string) Model {
	return Model{help: help}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return "📚 Help\n\n" + m.help()
}
