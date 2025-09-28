package credentials

import (
	"gophkeeper/client/internal/tui/common"

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
				input: common.InputWithPlaceholder("username"),
			},
			{
				title: "Password",
				input: func() textinput.Model {
					input := common.InputWithPlaceholder("username")
					input.EchoMode = textinput.EchoPassword
					input.EchoCharacter = '•'
					return input
				}(),
			},
			{
				title: "Note",
				input: common.InputWithPlaceholder("Any additional information"),
			},
		},
	}
}

func (m Model) Init() tea.Cmd {
	return m.fields[0].input.Focus()
}
