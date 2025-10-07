package sync

import "fmt"

// View renders the sync screen based on the current state.
func (m Model) View() string {
	title := "📡 SYNCING"
	footer := "\n(use ←/esc to return)"
	switch m.state {
	case stateInProgress:
		return fmt.Sprintf("%s\n\n 📡 Syncing...\n%s", title, footer)
	case stateSuccess:
		return fmt.Sprintf("%s\n\n ✅ Success!\n%s", title, footer)
	case stateError:
		return fmt.Sprintf("%s\n\n ❌ Error: %v\n%s", title, m.err, footer)
	default:
		return fmt.Sprintf("%s\n\n Initializing...\n%s", title, footer)
	}
}
