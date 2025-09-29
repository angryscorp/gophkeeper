package auth

import "fmt"

func (m Model) View() string {
	switch m.state {
	case stateAskUsername, stateAskPassword:
		return fmt.Sprintf("%s\n\n%s", m.title, m.input.View())
	case stateInProgress:
		return fmt.Sprintf("%s\n\nsending request...", m.title)
	case stateSuccess:
		return fmt.Sprintf("%s\n\nRequest is successful", m.title)
	case stateError:
		return fmt.Sprintf("%s\n\nRequest failed: %v", m.title, m.err)
	default:
		return fmt.Sprintf("%s\n\ninitializing...", m.title)
	}
}
