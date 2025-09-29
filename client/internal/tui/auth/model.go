package auth

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	title    string
	action   func(username, password string) error
	state    state
	input    textinput.Model
	err      error
	username string
	password string
}

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

func (m Model) Init() tea.Cmd {
	return cmdAskUsername.Run
}
