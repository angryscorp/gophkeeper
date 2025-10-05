package auth

import "fmt"

func (m Model) View() string {
	footer := "\n(use enter to submit, ←/esc to return)"
	switch m.state {
	case stateAskUsername, stateAskPassword:
		return fmt.Sprintf("%s\n\n%s\n%s", m.title, m.input.View(), footer)
	case stateInProgress:
		return fmt.Sprintf("%s\n\n 📡 Sending request...\n%s", m.title, footer)
	case stateSuccess:
		return fmt.Sprintf("%s\n\n ✅ Success!\n%s", m.title, footer)
	case stateError:
		return fmt.Sprintf("%s\n\n ❌ Error: %v\n%s", m.title, m.err, footer)
	default:
		return fmt.Sprintf("%s\n\nInitializing...\n%s", m.title, footer)
	}
}
