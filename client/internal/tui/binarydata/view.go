package binarydata

import "fmt"

func (m Model) View() string {
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

	view := "📁 Binary Data\n\n" +
		"File Path:\n" + m.fields[0].input.View() + "\n"

	if fileStatus != "" {
		view += fileStatus + "\n"
	}

	view += "\nDescription:\n" + m.fields[1].input.View() + "\n\n"

	if m.fileExists {
		view += "(tab to switch fields, ctrl+s to save, ←/esc to return)"
	} else {
		view += "(enter valid file path to enable saving, tab to switch fields, ←/esc to return)"
	}

	return view
}
