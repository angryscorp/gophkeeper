package bankcard

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
	saver     func(domain.BankCard) error
	resultMsg string
}

func (m Model) Focused() int {
	return m.focused
}
func New(saver func(domain.BankCard) error) Model {
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
