package textdata

func (m Model) View() string {
	titleStyle := ""
	contentStyle := ""

	if m.focused == 0 {
		titleStyle = " (focused)"
	} else {
		contentStyle = " (focused)"
	}

	return "📝 Text Data\n\n" +
		"Title" + titleStyle + ":\n" + m.title.View() + "\n\n" +
		"Content" + contentStyle + ":\n" + m.content.View() + "\n\n" +
		"(tab to switch fields, ctrl+s to save, ←/esc to return)"
}
