package bankcard

import (
	"errors"

	"gophkeeper/client/internal/domain"
	"gophkeeper/client/internal/tui/common"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

// Update implements tea.Model. It handles navigation between fields (up/down),
// submit with Ctrl+S (validation + save), and routes other key events into the
// focused text input.
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

	bankCard := domain.BankCard{
		ID:         uuid.New(),
		Owner:      m.fields[0].input.Value(),
		Number:     m.fields[1].input.Value(),
		ExpireDate: m.fields[2].input.Value(),
		CVV:        m.fields[3].input.Value(),
		Note:       m.fields[4].input.Value(),
	}

	err := m.saver(bankCard)
	if err != nil {
		return err
	}

	return nil
}
