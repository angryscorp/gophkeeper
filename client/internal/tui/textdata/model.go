package textdata

import (
	"gophkeeper/client/internal/domain"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Model is a Bubble Tea model for creating and saving text records.
type Model struct {
	title     textinput.Model
	content   textarea.Model
	focused   int // 0 = title, 1 = content
	saver     func(domain.UserTextData) error
	resultMsg string
}

// New creates a textdata form with title and multi-line content inputs.
func New(saver func(domain.UserTextData) error) Model {
	title := textinput.New()
	title.Placeholder = "Enter title..."
	title.Width = len(title.Placeholder)
	title.Focus()

	content := textarea.New()
	content.Placeholder = "Enter your text content here...\n\nYou can write multiple lines."
	content.SetWidth(80)
	content.SetHeight(10)
	content.ShowLineNumbers = false

	return Model{
		title:   title,
		content: content,
		focused: 0,
		saver:   saver,
	}
}

// Init sets up initial commands (blinking cursors for inputs).
func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, textarea.Blink)
}
