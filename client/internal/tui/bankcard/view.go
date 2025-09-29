package bankcard

import (
	"strings"
)

func (m Model) View() string {
	var b strings.Builder
	b.WriteString("💳 Bank Card Information\n\n")

	for i, field := range m.fields {
		b.WriteString(field.title + ":\n")
		b.WriteString(field.input.View())
		if i < len(m.fields)-1 {
			b.WriteRune('\n')
		}
		b.WriteRune('\n')
	}

	b.WriteString("\n(tab/shift+tab to navigate, ctrl+s to save, ←/esc to return)")
	return b.String()
}
