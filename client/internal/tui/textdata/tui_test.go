package textdata

import (
	"strings"
	"testing"

	"gophkeeper/client/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInit_NoPanic(t *testing.T) {
	m := New(func(domain.UserTextData) error { return nil })
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Init panicked: %v", r)
		}
	}()
	_ = m.Init()
}

func TestValidation_TitleAndContent(t *testing.T) {
	var called int
	m := New(func(domain.UserTextData) error { called++; return nil })

	m = sendKey(m, tea.KeyCtrlS)
	if out := m.View(); !strings.Contains(out, "Error saving data: title cannot be empty") {
		t.Fatalf("want title validation error, got:\n%s", out)
	}
	if called != 0 {
		t.Fatalf("saver must not be called on title validation error")
	}

	m.title.SetValue("My title")
	m = sendKey(m, tea.KeyCtrlS)
	if out := m.View(); !strings.Contains(out, "Error saving data: content cannot be empty") {
		t.Fatalf("want content validation error, got:\n%s", out)
	}
	if called != 0 {
		t.Fatalf("saver must not be called on content validation error")
	}
}

func TestSaveSuccessAndReset(t *testing.T) {
	var (
		called int
		got    domain.UserTextData
	)
	m := New(func(u domain.UserTextData) error { called++; got = u; return nil })

	m.title.SetValue("Title X")
	m.content.SetValue("Hello\nWorld")

	// Saving
	m = sendKey(m, tea.KeyCtrlS)

	if called != 1 {
		t.Fatalf("saver must be called once, got %d", called)
	}
	if got.Note != "Title X" || got.Data != "Hello\nWorld" {
		t.Fatalf("saver received wrong payload: %+v", got)
	}
	if out := m.View(); !strings.Contains(out, "Data successfully saved") {
		t.Fatalf("want success message, got:\n%s", out)
	}

	// Enter should reset the message
	m = sendKey(m, tea.KeyEnter)
	if out := m.View(); strings.Contains(out, "Data successfully saved") {
		t.Fatalf("success message must be cleared after Enter")
	}

	// All fields are reset
	if m.title.Value() != "" || m.content.Value() != "" {
		t.Fatalf("inputs must be reset, got title=%q content=%q", m.title.Value(), m.content.Value())
	}
}

func TestFocusSwitching(t *testing.T) {
	m := New(func(domain.UserTextData) error { return nil })

	// Focus
	if out := m.View(); !strings.Contains(out, "Title (focused):") {
		t.Fatalf("want title focused initially, got:\n%s", out)
	}

	// Tab -> change focus
	m = sendKey(m, tea.KeyTab)
	if out := m.View(); !strings.Contains(out, "Content (focused):") {
		t.Fatalf("want content focused after Tab, got:\n%s", out)
	}

	// Shift+Tab -> focus back
	m = sendKey(m, tea.KeyShiftTab)
	if out := m.View(); !strings.Contains(out, "Title (focused):") {
		t.Fatalf("want title focused after Shift+Tab, got:\n%s", out)
	}
}

func TestHeaderPresentInView(t *testing.T) {
	m := New(func(domain.UserTextData) error { return nil })
	if out := m.View(); !strings.HasPrefix(out, "📝 Text Data") {
		t.Fatalf("header missing, got:\n%s", out)
	}
}

func sendKey(m Model, k tea.KeyType) Model {
	nm, _ := m.Update(tea.KeyMsg{Type: k})
	return nm.(Model)
}
