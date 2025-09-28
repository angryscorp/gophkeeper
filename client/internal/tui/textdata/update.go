package textdata

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			// Switch focus between title and content
			if m.focused == 0 {
				m.focused = 1
				m.title.Blur()
				cmds = append(cmds, m.content.Focus())
			} else {
				m.focused = 0
				m.content.Blur()
				cmds = append(cmds, m.title.Focus())
			}
		}
	}

	// Update the focused model
	var cmd tea.Cmd
	if m.focused == 0 {
		m.title, cmd = m.title.Update(msg)
	} else {
		m.content, cmd = m.content.Update(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
