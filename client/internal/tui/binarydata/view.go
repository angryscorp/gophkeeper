package binarydata

import (
	"fmt"
	"strings"
)

// View implements tea.Model. It renders the binary data form UI,
// showing the file path input, file existence/size status, description
// input, and context-sensitive instructions at the bottom.
func (m Model) View() string {
	var b strings.Builder
	b.WriteString("📁 Binary Data\n\n")

	if m.resultMsg != "" {
		b.WriteString(m.resultMsg)
		b.WriteString("\n\n(use enter/space to reset, ←/esc to return)")
		return b.String()
	}

	b.WriteString("File Path:\n" + m.fields[0].input.View() + "\n")

	fileStatus := ""
	filePath := m.fields[0].input.Value()

	if filePath != "" {
		if m.fileExists {
			fileStatus = " ✅ File found"
			if m.fileSize > 0 {
				fileStatus += fmt.Sprintf(" (%d B)", m.fileSize)
			}
		} else {
			fileStatus = " ❌ File not found"
		}
	}

	if fileStatus != "" {
		b.WriteString(fileStatus + "\n")
	}

	b.WriteString("\nDescription:\n" + m.fields[1].input.View() + "\n\n")

	if m.fileExists {
		b.WriteString("(tab to switch fields, ctrl+s to save, ←/esc to return)")
	} else {
		b.WriteString("(enter valid file path to enable saving, tab to switch fields, ←/esc to return)")
	}

	return b.String()
}
