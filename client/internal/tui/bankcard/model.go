package bankcard

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

func (m Model) Focused() int {
	return m.focused
}
func New() Model {
	return Model{
		fields: []field{
			{
				title: "Cardholder Name",
				input: common.InputWithPlaceholder("John Doe"),
			},
			{
				title: "Card Number",
				input: common.InputWithPlaceholder("1234 5678 9012 3456"),
			},
			{
				title: "Expiry",
				input: common.InputWithPlaceholder("MM/YY"),
			},
			{
				title: "CVC",
				input: common.InputWithPlaceholder("123"),
			},
		},
	}
}

func (m Model) Init() tea.Cmd {
	return m.fields[0].input.Focus()
}
