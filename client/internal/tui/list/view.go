package list

import (
	"fmt"
)

// View implements tea.Model. It renders the list screen UI
// depending on the current state: loading in progress,
// successful data display, or an error message.
func (m Model) View() string {
	title := "📒 Private Data"
	footer := "\n(use ←/esc to return)"
	switch m.state {
	case stateInProgress:
		return fmt.Sprintf("%s\n\n 📡 Getting data...\n%s", title, footer)
	case stateSuccess:
		return m.list.View()
	case stateError:
		return fmt.Sprintf("%s\n\n ❌ Error: %v\n%s", title, m.err, footer)
	default:
		return fmt.Sprintf("%s\n\n Initializing...\n%s", title, footer)
	}
}
