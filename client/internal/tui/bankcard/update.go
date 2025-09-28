package bankcard

import (
	tea "github.com/charmbracelet/bubbletea"
)

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
		case "down", "enter":
			if m.focused < len(m.fields)-1 {
				m.fields[m.focused].input.Blur()
				m.focused++
				cmd := m.fields[m.focused].input.Focus()
				return m, cmd
			}
		}
	}

	var cmd tea.Cmd
	m.fields[m.focused].input, cmd = m.fields[m.focused].input.Update(msg)
	if m.fields[m.focused].input.Width < len(m.fields[m.focused].input.Value()) {
		m.fields[m.focused].input.Width = len(m.fields[m.focused].input.Value())
	}
	return m, cmd
}
