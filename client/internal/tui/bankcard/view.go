package bankcard

import (
	"strings"
)

// View implements tea.Model. It renders the bank card form UI.
func (m Model) View() string {
	var b strings.Builder
	b.WriteString("💳 Bank Card Information\n\n")

	if m.resultMsg != "" {
		b.WriteString(m.resultMsg)
		b.WriteString("\n\n(use enter/space to reset, ←/esc to return)")
		return b.String()
	}

	for i, field := range m.fields {
		b.WriteString(field.title + ":\n")
		b.WriteString(field.input.View())
		if i < len(m.fields)-1 {
			b.WriteRune('\n')
		}
		b.WriteRune('\n')
	}

	b.WriteString("\n(use ↑/↓ to navigate, ctrl+s to save, ←/esc to return)")
	return b.String()
}
