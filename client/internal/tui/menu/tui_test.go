package menu

import (
	"strings"
	"testing"

	"gophkeeper/client/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
)

type dummySaver struct{}

func (d dummySaver) SaveBankCard(_ domain.BankCard) error             { return nil }
func (d dummySaver) SaveCredentials(_ domain.Credentials) error       { return nil }
func (d dummySaver) SaveUserBinaryData(_ domain.UserBinaryData) error { return nil }
func (d dummySaver) SaveUserTextData(_ domain.UserTextData) error     { return nil }

func makeEnv() Environment {
	return Environment{
		RegFactory:   func(username, password string) error { return nil },
		LoginFactory: func(username, password string) error { return nil },
		DataSaver:    dummySaver{},
		SyncFactory:  func() error { return nil },
		DataFactory:  func() ([]domain.Record, error) { return nil, nil },
		HelpFactory:  func() string { return "Test help content" },
	}
}

func TestNewMenu_Defaults(t *testing.T) {
	m := New(makeEnv())

	if m.route != routeMenu {
		t.Fatalf("route = %v, want %v", m.route, routeMenu)
	}
	if m.cursor != 1 {
		t.Fatalf("cursor = %d, want 1", m.cursor)
	}
	if len(m.items) < 5 {
		t.Fatalf("len(items) = %d, want >= 5", len(m.items))
	}
}

func TestCursorBounds(t *testing.T) {
	m := New(makeEnv())

	// pressing up at top shouldn't go below 1
	m = sendKey(m, tea.KeyUp)
	if m.cursor != 1 {
		t.Fatalf("cursor after up at top = %d, want 1", m.cursor)
	}

	// press down many times, should clamp at len(items)
	for i := 0; i < 100; i++ {
		m = sendKey(m, tea.KeyDown)
	}
	if m.cursor != len(m.items) {
		t.Fatalf("cursor after many downs = %d, want %d", m.cursor, len(m.items))
	}
}

func TestEnterRegisterTransitions(t *testing.T) {
	m := New(makeEnv())

	// First item is Register in current menu definition
	if m.items[0].route != routeRegister {
		t.Fatalf("first item route = %v, want routeRegister", m.items[0].route)
	}

	// Ensure cursor points to first item, then Enter
	m.cursor = 1
	m = sendKey(m, tea.KeyEnter)

	if m.route != routeRegister {
		t.Fatalf("route after enter = %v, want %v", m.route, routeRegister)
	}
}

func TestHelpRenders(t *testing.T) {
	m := New(makeEnv())

	// Find the Help item index dynamically (order might change later)
	idx := -1
	for i, it := range m.items {
		if it.route == routeHelp {
			idx = i + 1 // cursor is 1-based
			break
		}
	}
	if idx == -1 {
		t.Skip("Help item not found; skipping")
	}

	m.cursor = idx
	m = sendKey(m, tea.KeyEnter) // enter Help

	if m.route != routeHelp {
		t.Fatalf("route after enter help = %v, want %v", m.route, routeHelp)
	}

	out := m.View()
	if !strings.Contains(out, "📚 Help") && !strings.Contains(out, "Help") {
		t.Fatalf("help view doesn't look like help; got:\n%s", out)
	}
}

// helper to pass a key to Update
func sendKey(m Model, key tea.KeyType) Model {
	next, _ := m.Update(tea.KeyMsg{Type: key})
	return next.(Model)
}
