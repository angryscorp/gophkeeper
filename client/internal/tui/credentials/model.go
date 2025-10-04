package credentials

import (
	"gophkeeper/client/internal/domain"
	"gophkeeper/client/internal/tui/common"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type field struct {
	title string
	input textinput.Model
}

type Model struct {
	fields    []field
	focused   int
	saver     func(domain.Credentials) error
	resultMsg string
}

func New(saver func(domain.Credentials) error) Model {
	return Model{
		fields: []field{
			{
				title: "Username",
				input: common.InputWithPlaceholder("username"),
			},
			{
				title: "Password",
				input: func() textinput.Model {
					input := common.InputWithPlaceholder("password")
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
		saver: saver,
	}
}

func (m Model) Init() tea.Cmd {
	return m.fields[0].input.Focus()
}
