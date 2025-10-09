package list

import (
	"errors"
	"strings"
	"testing"

	"gophkeeper/client/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInit_ReturnsCmd(t *testing.T) {
	m := New(func() ([]domain.Record, error) { return nil, nil })
	if cmd := m.Init(); cmd == nil {
		t.Fatalf("Init() returned nil command")
	}
}

func TestView_ByState(t *testing.T) {
	cases := []struct {
		name    string
		state   state
		wantSub string
	}{
		{"init", stateInit, "Initializing"},
		{"in-progress", stateInProgress, "Getting data"},
		{"success", stateSuccess, ""},
		{"error", stateError, "Error"},
	}

	m := New(func() ([]domain.Record, error) { return nil, nil })

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m.state = tc.state
			out := m.View()
			if !strings.Contains(out, tc.wantSub) {
				t.Fatalf("View() missing substring %q for state %v;\noutput:\n%s", tc.wantSub, tc.state, out)
			}
		})
	}
}

func TestUpdate_OnCmdLoadData_SetsInProgress(t *testing.T) {
	m := New(func() ([]domain.Record, error) { return nil, nil })

	var cmd tea.Cmd
	updated, cmd := m.Update(cmdLoadData)
	m = updated.(Model)

	if cmd == nil {
		t.Fatalf("expected a follow-up cmd after cmdLoadData")
	}
	if m.state != stateInProgress {
		t.Fatalf("state = %v, want %v", m.state, stateInProgress)
	}
}

func TestUpdate_OnResultMsg_Success(t *testing.T) {
	m := New(func() ([]domain.Record, error) { return nil, nil })

	msg := resultMsg{data: nil, err: nil}
	updated, _ := m.Update(msg)
	m = updated.(Model)

	if m.state != stateSuccess {
		t.Fatalf("state = %v, want %v", m.state, stateSuccess)
	}
}

func TestUpdate_OnResultMsg_Error(t *testing.T) {
	m := New(func() ([]domain.Record, error) { return nil, nil })

	msg := resultMsg{data: nil, err: errors.New("boom")}
	updated, _ := m.Update(msg)
	m = updated.(Model)

	if m.state != stateError {
		t.Fatalf("state = %v, want %v", m.state, stateError)
	}
}
