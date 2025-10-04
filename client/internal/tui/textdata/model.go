package textdata

import (
	"gophkeeper/client/internal/domain"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	title     textinput.Model
	content   textarea.Model
	focused   int // 0 = title, 1 = content
	saver     func(domain.UserTextData) error
	resultMsg string
}

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

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, textarea.Blink)
}
