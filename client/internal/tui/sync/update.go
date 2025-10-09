package sync

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles sync commands and results, updating the state
// to InProgress, Success, or Error, and returning the next command.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case baseCmd:
		switch msg {
		case cmdDoSync:
			m.state = stateInProgress
			return m, func() tea.Msg {
				err := m.sync()
				return resultMsg{
					success: err == nil,
					err:     err,
				}
			}
		}

	case resultMsg:
		if msg.success {
			m.state = stateSuccess
		} else {
			m.err = msg.err
			m.state = stateError
		}
		return m, nil
	}

	return m, nil
}
