package list

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// baseCmd represents internal commands used by the list model.
// It drives async actions like loading data.
type baseCmd int

const (
	cmdLoadData baseCmd = iota
)

// Run implements tea.Cmd for baseCmd by returning the command
// as a tea.Msg, enabling it to be processed in Update.
func (m baseCmd) Run() tea.Msg {
	return m
}

// resultMsg is a Bubble Tea message carrying the result of a
// data-loading operation: a slice of list items or an error.
type resultMsg struct {
	data []list.Item
	err  error
}

// Run implements tea.Cmd for resultMsg by returning itself
// as a tea.Msg so it can be dispatched in the update loop.
func (c resultMsg) Run() tea.Msg {
	return c
}
