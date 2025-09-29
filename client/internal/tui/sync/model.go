package sync

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	sync  func() error
	state state
	err   error
}

func New(sync func() error) Model {
	return Model{sync: sync}
}

func (m Model) Init() tea.Cmd {
	return cmdDoSync.Run
}
