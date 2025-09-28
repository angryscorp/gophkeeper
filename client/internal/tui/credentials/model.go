package credentials

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type field struct {
	title string
	input textinput.Model
}

type Model struct {
	fields  []field
	focused int
}

func New() Model {

	return Model{
		fields: []field{
			{
				title: "Username",
				input: initInputWithPlaceholder("username"),
			},
			{
				title: "Password",
				input: func() textinput.Model {
					input := initInputWithPlaceholder("username")
					input.EchoMode = textinput.EchoPassword
					input.EchoCharacter = '•'
					return input
				}(),
			},
			{
				title: "Note",
				input: initInputWithPlaceholder("Any additional information"),
			},
		},
	}
}

func (m Model) Init() tea.Cmd {
	return m.fields[0].input.Focus()
}

func initInputWithPlaceholder(placeholder string) textinput.Model {
	input := textinput.New()
	input.Placeholder = placeholder
	input.Width = len(placeholder)
	return input
}
