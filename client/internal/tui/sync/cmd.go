package sync

import tea "github.com/charmbracelet/bubbletea"

// baseCmd represents internal commands for the sync model.
type baseCmd int

const (
	cmdDoSync baseCmd = iota
)

// Run implements tea.Cmd for baseCmd by returning itself as a tea.Msg.
func (m baseCmd) Run() tea.Msg {
	return m
}

// resultMsg carries the outcome of a sync operation: either success or error.
type resultMsg struct {
	success bool
	err     error
}

// Run implements tea.Cmd for resultMsg by returning itself as a tea.Msg.
func (c resultMsg) Run() tea.Msg {
	return c
}
