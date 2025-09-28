package binarydata

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.fields[m.focused].input.Blur()
			newVal := m.focused + 1
			if newVal == 2 {
				newVal = 0
			}
			m.focused = newVal
			cmd := m.fields[m.focused].input.Focus()
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.fields[m.focused].input, cmd = m.fields[m.focused].input.Update(msg)
	currentWidth := len(m.fields[m.focused].input.Value())
	if m.fields[m.focused].input.Width < currentWidth {
		m.fields[m.focused].input.Width = currentWidth
	}

	// Check file existence when a file path changes
	if m.focused == 0 {
		m.fileExists, m.fileSize = checkFilePath(m.fields[0].input.Value())
	}

	return m, cmd
}

func checkFilePath(filePath string) (fileExists bool, fileSize int64) {
	if filePath == "" {
		return
	}
	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		return
	}
	return true, info.Size()
}
