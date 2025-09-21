package auth

import tea "github.com/charmbracelet/bubbletea"

type baseCmd int

const (
	cmdAskUsername baseCmd = iota
	cmdSendRequest
)

func (m baseCmd) Run() tea.Msg {
	return m
}

type resultMsg struct {
	success bool
	err     error
}

func (c resultMsg) Run() tea.Msg {
	return c
}
