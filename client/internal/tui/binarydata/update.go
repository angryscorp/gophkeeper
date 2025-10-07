package binarydata

import (
	"errors"
	"gophkeeper/client/internal/domain"
	"gophkeeper/client/internal/tui/common"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

// Update implements tea.Model. It handles navigation (tab/shift+tab),
// saving with Ctrl+S, and routes other events to the focused input.
// Also updates file existence/size when the file path changes.
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
		case "ctrl+s":
			err := m.saveData()
			if err != nil {
				m.resultMsg = "Error saving data: " + err.Error()
				return m, nil
			}
			m.resultMsg = "Data successfully saved"

			// reset all fields
			for i := 0; i < len(m.fields); i++ {
				m.fields[i].input.Reset()
				m.fields[m.focused].input.Blur()
			}
			m.focused = 0
			return m, m.fields[m.focused].input.Focus()

		case "enter", " ":
			m.resultMsg = ""
			return m, nil
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

func (m Model) saveData() error {
	for _, f := range m.fields {
		if f.input.Value() == "" {
			return errors.New(f.title + " is empty")
		}
	}

	data, err := os.ReadFile(m.fields[0].input.Value())
	if err != nil {
		return err
	}

	binaryData := domain.UserBinaryData{
		ID:   uuid.New(),
		Data: data,
		Note: m.fields[1].input.Value(),
	}

	err = m.saver(binaryData)
	if err != nil {
		return err
	}

	return nil
}
