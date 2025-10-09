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

// Model is a Bubble Tea model that collects a file path and a note,
// tracks whether the file exists and its size, and saves via the
// supplied saver callback.
type Model struct {
	fields     []field
	focused    int
	fileExists bool
	fileSize   int64
	saver      func(domain.UserBinaryData) error
	resultMsg  string
}

// New creates a new binary-data form model with sensible placeholders
// and the saver callback to persist the collected data.
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

// Init implements tea.Model. It sets initial focus to the first field.
func (m Model) Init() tea.Cmd {
	return m.fields[0].input.Focus()
}
