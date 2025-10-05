package list

import (
	"fmt"
)

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
