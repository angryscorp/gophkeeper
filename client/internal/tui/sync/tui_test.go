package sync

import (
	"errors"
	"strings"
	"testing"
)

func TestInit_ReturnsCommand(t *testing.T) {
	m := Model{sync: func() error { return nil }}
	if cmd := m.Init(); cmd == nil {
		t.Fatalf("Init() returned nil cmd; want non-nil")
	}
}

func TestUpdate_SuccessFlow(t *testing.T) {
	m := Model{sync: func() error { return nil }}

	// Init returns a cmd that emits baseCmd (cmdDoSync)
	cmd := m.Init()
	if cmd == nil {
		t.Fatalf("Init() returned nil cmd")
	}

	// Feed that message into Update (this switches to InProgress and returns a cmd that runs sync)
	msg := cmd() // baseCmd
	nm, nextCmd := m.Update(msg)
	m2, ok := nm.(Model)
	if !ok {
		t.Fatalf("Update() did not return Model type")
	}
	if m2.state != stateInProgress {
		t.Fatalf("state after first Update = %v; want stateInProgress", m2.state)
	}

	// Execute returned cmd (runs m.sync) and feed resultMsg back into Update
	if nextCmd == nil {
		t.Fatalf("expected non-nil cmd from Update after cmdDoSync")
	}
	resMsg := nextCmd()
	nm, _ = m2.Update(resMsg)
	m3, ok := nm.(Model)
	if !ok {
		t.Fatalf("Update() did not return Model type on resultMsg")
	}
	if m3.state != stateSuccess {
		t.Fatalf("final state = %v; want stateSuccess", m3.state)
	}

	// View should reflect success
	out := m3.View()
	if !strings.Contains(out, "Success") && !strings.Contains(out, "✅") {
		t.Fatalf("View() = %q; want to contain success indicator", out)
	}
}

func TestUpdate_ErrorFlow(t *testing.T) {
	m := Model{sync: func() error { return errors.New("boom") }}

	// Kick off
	cmd := m.Init()
	msg := cmd()

	nm, nextCmd := m.Update(msg)
	m2, ok := nm.(Model)
	if !ok {
		t.Fatalf("Update() did not return Model type")
	}
	if m2.state != stateInProgress {
		t.Fatalf("state after first Update = %v; want stateInProgress", m2.state)
	}

	// Run the returned cmd (produces resultMsg with error)
	resMsg := nextCmd()
	nm, _ = m2.Update(resMsg)
	m3, ok := nm.(Model)
	if !ok {
		t.Fatalf("Update() did not return Model type on resultMsg")
	}
	if m3.state != stateError {
		t.Fatalf("final state = %v; want stateError", m3.state)
	}
	if m3.err == nil {
		t.Fatalf("expected non-nil error stored in model")
	}

	// View should reflect error
	out := m3.View()
	if !strings.Contains(out, "Error") && !strings.Contains(out, "❌") {
		t.Fatalf("View() = %q; want to contain error indicator", out)
	}
}

func TestView_ByState(t *testing.T) {
	tests := []struct {
		name   string
		model  Model
		expect []string
	}{
		{
			name:   "init",
			model:  Model{}, // zero state
			expect: []string{"Initializing", "📡"},
		},
		{
			name:   "in-progress",
			model:  Model{state: stateInProgress},
			expect: []string{"Syncing", "📡"},
		},
		{
			name:   "success",
			model:  Model{state: stateSuccess},
			expect: []string{"Success", "✅"},
		},
		{
			name:   "error",
			model:  Model{state: stateError, err: errors.New("x")},
			expect: []string{"Error", "❌"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := tc.model.View()
			for _, sub := range tc.expect {
				if !strings.Contains(out, sub) {
					t.Fatalf("View() missing substring %q;\noutput:\n%s", sub, out)
				}
			}
		})
	}
}

func TestNoPanic_OnUnknownMsg(t *testing.T) {
	m := Model{sync: func() error { return nil }}
	// Send some unrelated message; Update should be safe
	_, _ = m.Update(struct{ foo string }{foo: "bar"})
}
