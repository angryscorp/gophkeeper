package binarydata

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type field struct {
	title string
	input textinput.Model
}

type Model struct {
	fields     []field
	focused    int
	fileExists bool
	fileSize   int64
}

func New() Model {
	return Model{
		fields: []field{
			{
				title: "File Path",
				input: initInputWithPlaceholder("/path/to/file"),
			},
			{
				title: "Note",
				input: initInputWithPlaceholder("File description"),
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
