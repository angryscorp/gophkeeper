package help

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Model is a Bubble Tea model that renders a static help screen.
// It delegates the actual help text generation to the provided
// callback function.
type Model struct {
	help func() string
}

// New creates a new help Model using the given function to
// generate the help text.
func New(help func() string) Model {
	return Model{help: help}
}

// Init implements tea.Model. It performs no initialization.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model. It does nothing and returns the
// model unchanged, since the help screen is static.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// View implements tea.Model. It renders the help screen with a
// header and the provided help text.
func (m Model) View() string {
	return "📚 Help\n\n" + m.help()
}
