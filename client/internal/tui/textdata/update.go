package textdata

import (
	"errors"

	"gophkeeper/client/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

// Update handles keyboard input and routes it to the focused field (title or content).
// It also supports focus switching (tab/shift+tab) and saving (ctrl+s).
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			// Switch focus between title and content
			if m.focused == 0 {
				m.focused = 1
				m.title.Blur()
				cmds = append(cmds, m.content.Focus())
			} else {
				m.focused = 0
				m.content.Blur()
				cmds = append(cmds, m.title.Focus())
			}
		case "ctrl+s":
			err := m.saveData()
			if err != nil {
				m.resultMsg = "Error saving data: " + err.Error()
				return m, nil
			}
			m.resultMsg = "Data successfully saved"

			// reset all fields
			m.title.Reset()
			m.content.Reset()
			m.content.Blur()
			m.focused = 0

			return m, m.title.Focus()

		case "enter", " ":
			if m.resultMsg != "" {
				m.resultMsg = ""
				return m, nil
			}
		}
	}

	// Update the focused model
	var cmd tea.Cmd
	if m.focused == 0 {
		m.title, cmd = m.title.Update(msg)
	} else {
		m.content, cmd = m.content.Update(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) saveData() error {
	if m.title.Value() == "" {
		return errors.New("title cannot be empty")
	}

	if m.content.Value() == "" {
		return errors.New("content cannot be empty")
	}

	bankCard := domain.UserTextData{
		ID:   uuid.New(),
		Data: m.content.Value(),
		Note: m.title.Value(),
	}

	err := m.saver(bankCard)
	if err != nil {
		return err
	}

	return nil
}
