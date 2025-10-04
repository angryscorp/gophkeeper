package sync

import "fmt"

func (m Model) View() string {
	title := "📡 SYNCING"
	switch m.state {
	case stateInProgress:
		return fmt.Sprintf("%s\n\nsyncing...", title)
	case stateSuccess:
		return fmt.Sprintf("%s\n\nSuccess", title)
	case stateError:
		return fmt.Sprintf("%s\n\nFailure: %v", title, m.err)
	default:
		return fmt.Sprintf("%s\n\ninitializing...", title)
	}
}
