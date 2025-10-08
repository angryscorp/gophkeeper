package bankcard

import (
	"errors"
	"gophkeeper/client/internal/domain"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModel_Init(t *testing.T) {
	m := New(func(bc domain.BankCard) error { return nil })
	cmd := m.Init()

	if cmd == nil {
		t.Fatal("Init() returned nil cmd")
	}
}

func TestModel_Update_Navigation(t *testing.T) {
	m := New(func(bc domain.BankCard) error { return nil })

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
			m := New(func(bc domain.BankCard) error { return tt.saverErr })

			if tt.fillFields {
				for i := range m.fields {
					m.fields[i].input.SetValue("test" + string(rune('0'+i)))
				}
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
	var savedCard domain.BankCard
	m := New(func(bc domain.BankCard) error {
		savedCard = bc
		return nil
	})

	// Fill all fields
	m.fields[0].input.SetValue("John Doe")
	m.fields[1].input.SetValue("1234567890123456")
	m.fields[2].input.SetValue("12/25")
	m.fields[3].input.SetValue("123")
	m.fields[4].input.SetValue("test note")

	err := m.saveData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if savedCard.Owner != "John Doe" {
		t.Errorf("got owner %q, want %q", savedCard.Owner, "John Doe")
	}
	if savedCard.Number != "1234567890123456" {
		t.Errorf("got number %q, want %q", savedCard.Number, "1234567890123456")
	}
}

func TestModel_View(t *testing.T) {
	tests := []struct {
		name      string
		resultMsg string
		wantText  string
	}{
		{"normal view", "", "Bank Card Information"},
		{"result view", "Success!", "Success!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(func(bc domain.BankCard) error { return nil })
			m.resultMsg = tt.resultMsg

			view := m.View()
			if !strings.Contains(view, tt.wantText) {
				t.Errorf("view doesn't contain %q", tt.wantText)
			}
		})
	}
}
