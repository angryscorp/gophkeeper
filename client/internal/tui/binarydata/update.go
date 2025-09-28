package binarydata

import (
	"gophkeeper/client/internal/tui/common"
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
	common.ExtendInputWidthIfNeeded(&m.fields[m.focused].input)

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
