package textdata

import "strings"

// View renders the text data form UI, showing title and content fields,
// highlighting the focused one and displaying result messages if present.
func (m Model) View() string {
	var b strings.Builder
	b.WriteString("📝 Text Data\n\n")

	if m.resultMsg != "" {
		b.WriteString(m.resultMsg)
		b.WriteString("\n\n(use enter/space to reset, ←/esc to return)")
		return b.String()
	}

	titleStyle := ""
	contentStyle := ""

	if m.focused == 0 {
		titleStyle = " (focused)"
	} else {
		contentStyle = " (focused)"
	}

	b.WriteString("Title" + titleStyle + ":\n" + m.title.View() + "\n\n")
	b.WriteString("Content" + contentStyle + ":\n" + m.content.View() + "\n\n")
	b.WriteString("(tab to switch fields, ctrl+s to save, ←/esc to return)")

	return b.String()
}
