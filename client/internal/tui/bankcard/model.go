package bankcard

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

func (m Model) Focused() int {
	return m.focused
}
func New() Model {
	return Model{
		fields: []field{
			{
				title: "Cardholder Name",
				input: initInputWithPlaceholder("John Doe"),
			},
			{
				title: "Card Number",
				input: initInputWithPlaceholder("1234 5678 9012 3456"),
			},
			{
				title: "Expiry",
				input: initInputWithPlaceholder("MM/YY"),
			},
			{
				title: "CVC",
				input: initInputWithPlaceholder("123"),
			},
		},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.fields[0].input.Focus(), textinput.Blink)
}

func initInputWithPlaceholder(placeholder string) textinput.Model {
	input := textinput.New()
	input.Placeholder = placeholder
	input.Width = len(placeholder)
	return input
}
