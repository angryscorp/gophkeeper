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

// Model is a Bubble Tea model for collecting bank card details
// from the user in a multi-field TUI form.
type Model struct {
	fields    []field
	focused   int
	saver     func(domain.BankCard) error
	resultMsg string
}

// New creates a new bank card form model with default placeholders
// and a saver callback to persist the collected card data.
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

// Init implements tea.Model. It sets the initial focus to the first field.
func (m Model) Init() tea.Cmd {
	return m.fields[0].input.Focus()
}
