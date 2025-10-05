package list

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type baseCmd int

const (
	cmdLoadData baseCmd = iota
)

func (m baseCmd) Run() tea.Msg {
	return m
}

type resultMsg struct {
	data []list.Item
	err  error
}

func (c resultMsg) Run() tea.Msg {
	return c
}
