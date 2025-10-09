package credentials

import (
	"errors"
	"strings"
	"testing"

	"gophkeeper/client/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModel_Init(t *testing.T) {
	m := New(func(c domain.Credentials) error { return nil })
	cmd := m.Init()

	if cmd == nil {
		t.Fatal("Init() returned nil cmd")
	}
}

func TestModel_Update_Navigation(t *testing.T) {
	m := New(func(c domain.Credentials) error { return nil })

	// Navigate down
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(Model)
	if m.focused != 1 {
		t.Errorf("got focused %d, want 1", m.focused)
	}

	// Navigate up
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(Model)
	if m.focused != 0 {
		t.Errorf("got focused %d, want 0", m.focused)
	}

	// Try up at boundary (should stay at 0)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(Model)
	if m.focused != 0 {
		t.Errorf("got focused %d, want 0", m.focused)
	}
}

func TestModel_Update_Save(t *testing.T) {
	tests := []struct {
		name       string
		saverErr   error
		fillFields bool
		wantMsg    string
	}{
		{"success", nil, true, "successfully saved"},
		{"empty fields", nil, false, "is empty"},
		{"saver error", errors.New("db error"), true, "Error saving"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(func(c domain.Credentials) error { return tt.saverErr })

			if tt.fillFields {
				m.fields[0].input.SetValue("user123")
				m.fields[1].input.SetValue("pass123")
				m.fields[2].input.SetValue("note")
			}

			updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			m = updated.(Model)

			if !strings.Contains(m.resultMsg, tt.wantMsg) {
				t.Errorf("resultMsg %q doesn't contain %q", m.resultMsg, tt.wantMsg)
			}
		})
	}
}

func TestModel_saveData(t *testing.T) {
	var savedCreds domain.Credentials
	m := New(func(c domain.Credentials) error {
		savedCreds = c
		return nil
	})

	m.fields[0].input.SetValue("testuser")
	m.fields[1].input.SetValue("testpass")
	m.fields[2].input.SetValue("test note")

	err := m.saveData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if savedCreds.Login != "testuser" {
		t.Errorf("got login %q, want %q", savedCreds.Login, "testuser")
	}
	if savedCreds.Password != "testpass" {
		t.Errorf("got password %q, want %q", savedCreds.Password, "testpass")
	}
	if savedCreds.Note != "test note" {
		t.Errorf("got note %q, want %q", savedCreds.Note, "test note")
	}
}

func TestModel_View(t *testing.T) {
	tests := []struct {
		name      string
		resultMsg string
		wantText  string
	}{
		{"normal view", "", "Credentials"},
		{"result view", "Saved!", "Saved!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(func(c domain.Credentials) error { return nil })
			m.resultMsg = tt.resultMsg

			view := m.View()
			if !strings.Contains(view, tt.wantText) {
				t.Errorf("view doesn't contain %q", tt.wantText)
			}
		})
	}
}
