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

// Model is a Bubble Tea form for entering user credentials
// (username, password, optional note) and saving them via
// the provided saver callback.
type Model struct {
	fields    []field
	focused   int
	saver     func(domain.Credentials) error
	resultMsg string
}

// New creates a credentials form Model with preconfigured
// input fields for username, password (masked), and note.
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

// Init implements tea.Model. It sets focus on the first input field.
func (m Model) Init() tea.Cmd {
	return m.fields[0].input.Focus()
}
