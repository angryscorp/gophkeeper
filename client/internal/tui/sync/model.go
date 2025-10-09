package sync

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Model is a Bubble Tea model for triggering and tracking
// a data synchronization operation with the server.
type Model struct {
	sync  func() error
	state state
	err   error
}

// New creates a new sync Model with the given function that
// performs the actual synchronization.
func New(sync func() error) Model {
	return Model{sync: sync}
}

// Init implements tea.Model. It starts the initial sync command.
func (m Model) Init() tea.Cmd {
	return cmdDoSync.Run
}
