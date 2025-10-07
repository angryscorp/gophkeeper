package auth

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the state machine for the authentication TUI flow.
// It manages username/password input fields, current state,
// error handling, and executes the provided action on submit.
type Model struct {
	title    string
	action   func(username, password string) error
	state    state
	input    textinput.Model
	err      error
	username string
	password string
}

// New creates a new authentication Model with the given title
// (e.g. "Login") and action callback. The action is called with
// username and password once the user completes input.
func New(
	title string,
	action func(username, password string) error,
) Model {
	return Model{
		title:  title,
		action: action,
		state:  stateInit,
		input:  textinput.New(),
	}
}

// Init is part of the Bubble Tea tea.Model interface.
// It returns the initial command, which starts by asking the username.
func (m Model) Init() tea.Cmd {
	return cmdAskUsername.Run
}
