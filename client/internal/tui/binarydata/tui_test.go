package binarydata

import (
	"errors"
	"gophkeeper/client/internal/domain"
	"os"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModel_Init(t *testing.T) {
	m := New(func(bd domain.UserBinaryData) error { return nil })
	cmd := m.Init()

	if cmd == nil {
		t.Fatal("Init() returned nil cmd")
	}
}

func TestModel_Update_Navigation(t *testing.T) {
	m := New(func(bd domain.UserBinaryData) error { return nil })

	// Navigate with tab
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)
	if m.focused != 1 {
		t.Errorf("got focused %d, want 1", m.focused)
	}

	// Tab again wraps to 0
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)
	if m.focused != 0 {
		t.Errorf("got focused %d, want 0", m.focused)
	}
}

func TestModel_Update_Save(t *testing.T) {
	// Create temp file for testing
	tmpFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("test content")
	tmpFile.Close()

	tests := []struct {
		name       string
		saverErr   error
		fillFields bool
		filePath   string
		wantMsg    string
	}{
		{"success", nil, true, tmpFile.Name(), "successfully saved"},
		{"empty fields", nil, false, "", "is empty"},
		{"saver error", errors.New("db error"), true, tmpFile.Name(), "Error saving"},
		{"file not found", nil, true, "/nonexistent/file.txt", "no such file"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(func(bd domain.UserBinaryData) error { return tt.saverErr })

			if tt.fillFields {
				m.fields[0].input.SetValue(tt.filePath)
				m.fields[1].input.SetValue("test note")
			}

			updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			m = updated.(Model)

			if !strings.Contains(m.resultMsg, tt.wantMsg) {
				t.Errorf("resultMsg %q doesn't contain %q", m.resultMsg, tt.wantMsg)
			}
		})
	}
}

func TestCheckFilePath(t *testing.T) {
	// Create temp file
	tmpFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("test")
	tmpFile.Close()

	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name       string
		path       string
		wantExists bool
		wantSize   int64
	}{
		{"valid file", tmpFile.Name(), true, 4},
		{"empty path", "", false, 0},
		{"nonexistent", "/nonexistent/file.txt", false, 0},
		{"directory", tmpDir, false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, size := checkFilePath(tt.path)

			if exists != tt.wantExists {
				t.Errorf("got exists %v, want %v", exists, tt.wantExists)
			}
			if size != tt.wantSize {
				t.Errorf("got size %d, want %d", size, tt.wantSize)
			}
		})
	}
}

func TestModel_saveData(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("test content")
	tmpFile.Close()

	var savedData domain.UserBinaryData
	m := New(func(bd domain.UserBinaryData) error {
		savedData = bd
		return nil
	})

	m.fields[0].input.SetValue(tmpFile.Name())
	m.fields[1].input.SetValue("test note")

	err = m.saveData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(savedData.Data) != "test content" {
		t.Errorf("got data %q, want %q", string(savedData.Data), "test content")
	}
	if savedData.Note != "test note" {
		t.Errorf("got note %q, want %q", savedData.Note, "test note")
	}
}

func TestModel_View(t *testing.T) {
	tests := []struct {
		name       string
		resultMsg  string
		fileExists bool
		wantText   string
	}{
		{"normal view", "", false, "Binary Data"},
		{"result view", "Saved!", false, "Saved!"},
		{"file exists", "", true, "File found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(func(bd domain.UserBinaryData) error { return nil })
			m.resultMsg = tt.resultMsg
			m.fileExists = tt.fileExists
			if tt.fileExists {
				m.fields[0].input.SetValue("test.txt")
			}

			view := m.View()
			if !strings.Contains(view, tt.wantText) {
				t.Errorf("view doesn't contain %q", tt.wantText)
			}
		})
	}
}
