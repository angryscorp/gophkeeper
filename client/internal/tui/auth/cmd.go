package auth

import tea "github.com/charmbracelet/bubbletea"

// baseCmd represents a command in the authentication state machine.
// It is implemented as an int enum and used to drive transitions
// between username input, password input, and request submission.
type baseCmd int

const (
	cmdAskUsername baseCmd = iota
	cmdAskPassword
	cmdSendRequest
)

// Run implements tea.Cmd for baseCmd by returning itself as a tea.Msg.
func (m baseCmd) Run() tea.Msg {
	return m
}

// resultMsg is a Bubble Tea message used to signal the outcome of an
// authentication attempt, carrying either a success flag or an error.
type resultMsg struct {
	success bool
	err     error
}

// Run implements tea.Cmd for resultMsg by returning itself as a tea.Msg.
func (c resultMsg) Run() tea.Msg {
	return c
}
