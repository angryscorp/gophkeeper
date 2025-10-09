package auth

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModel_Init(t *testing.T) {
	m := New("Test", func(u, p string) error { return nil })
	cmd := m.Init()
	msg := cmd()

	if msg != cmdAskUsername {
		t.Errorf("got msg %v, want cmdAskUsername", msg)
	}
}

func TestModel_Update_EmptyInputs(t *testing.T) {
	tests := []struct {
		name  string
		cmd   baseCmd
		input string
	}{
		{"empty username", cmdAskUsername, ""},
		{"empty password", cmdAskPassword, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New("Test", func(u, p string) error { return nil })
			updated, _ := m.Update(tt.cmd)
			m = updated.(Model)

			m.input.SetValue(tt.input)
			updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
			m = updated.(Model)

			if m.state != stateError {
				t.Errorf("got state %v, want stateError", m.state)
			}
			if m.err == nil {
				t.Error("expected error but got nil")
			}
		})
	}
}

func TestModel_Update_Result(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		wantState state
	}{
		{"success", nil, stateSuccess},
		{"failure", errors.New("failed"), stateError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New("Test", func(u, p string) error { return nil })
			msg := resultMsg{success: tt.err == nil, err: tt.err}
			updated, _ := m.Update(msg)
			m = updated.(Model)

			if m.state != tt.wantState {
				t.Errorf("got state %v, want %v", m.state, tt.wantState)
			}
		})
	}
}

func TestModel_View(t *testing.T) {
	tests := []struct {
		state    state
		wantText string
	}{
		{stateInit, "Initializing"},
		{stateInProgress, "Sending request"},
		{stateSuccess, "Success"},
		{stateError, "Error"},
	}

	for _, tt := range tests {
		m := Model{state: tt.state, err: errors.New("test")}
		view := m.View()

		if !strings.Contains(view, tt.wantText) {
			t.Errorf("view doesn't contain %q", tt.wantText)
		}
	}
}
