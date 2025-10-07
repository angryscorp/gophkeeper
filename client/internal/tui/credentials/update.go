package credentials

import (
	"errors"
	"gophkeeper/client/internal/domain"
	"gophkeeper/client/internal/tui/common"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

// Update implements tea.Model. It handles navigation between fields,
// saving data with Ctrl+S, clearing messages, and delegates input
// events to the currently focused field.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.focused > 0 {
				m.fields[m.focused].input.Blur()
				m.focused--
				cmd := m.fields[m.focused].input.Focus()
				return m, cmd
			}
		case "down":
			if m.focused < len(m.fields)-1 {
				m.fields[m.focused].input.Blur()
				m.focused++
				cmd := m.fields[m.focused].input.Focus()
				return m, cmd
			}
		case "ctrl+s":
			err := m.saveData()
			if err != nil {
				m.resultMsg = "Error saving data: " + err.Error()
				return m, nil
			}
			m.resultMsg = "Data successfully saved"

			// reset all fields
			for i := 0; i < len(m.fields); i++ {
				m.fields[i].input.Reset()
				m.fields[m.focused].input.Blur()
			}
			m.focused = 0
			return m, m.fields[m.focused].input.Focus()

		case "enter", " ":
			m.resultMsg = ""
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.fields[m.focused].input, cmd = m.fields[m.focused].input.Update(msg)
	common.ExtendInputWidthIfNeeded(&m.fields[m.focused].input)
	return m, cmd
}

func (m Model) saveData() error {
	for _, f := range m.fields {
		if f.input.Value() == "" {
			return errors.New(f.title + " is empty")
		}
	}

	credentials := domain.Credentials{
		ID:       uuid.New(),
		Login:    m.fields[0].input.Value(),
		Password: m.fields[1].input.Value(),
		Note:     m.fields[2].input.Value(),
	}

	err := m.saver(credentials)
	if err != nil {
		return err
	}

	return nil
}
