package binarydata

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
	fields     []field
	focused    int
	fileExists bool
	fileSize   int64
	saver      func(domain.UserBinaryData) error
	resultMsg  string
}

func New(saver func(domain.UserBinaryData) error) Model {
	return Model{
		fields: []field{
			{
				title: "File Path",
				input: common.InputWithPlaceholder("/path/to/file"),
			},
			{
				title: "Note",
				input: common.InputWithPlaceholder("File description"),
			},
		},
		saver: saver,
	}
}

func (m Model) Init() tea.Cmd {
	return m.fields[0].input.Focus()
}
