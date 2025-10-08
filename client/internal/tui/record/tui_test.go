package record

import (
	"strings"
	"testing"

	"gophkeeper/client/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
)

type stubSaver struct{}

func (stubSaver) SaveBankCard(_ domain.BankCard) error             { return nil }
func (stubSaver) SaveCredentials(_ domain.Credentials) error       { return nil }
func (stubSaver) SaveUserBinaryData(_ domain.UserBinaryData) error { return nil }
func (stubSaver) SaveUserTextData(_ domain.UserTextData) error     { return nil }

func TestInitialView_ShowsTypeSelectionHeader(t *testing.T) {
	m := New(stubSaver{})
	// Initial route should be the type selection menu; View should include the header.
	view := m.View()
	want := "Add New Item"
	if !strings.Contains(view, want) {
		t.Fatalf("initial View() should contain %q; got:\n%s", want, view)
	}
}

func TestEnterOnFirstItem_OpensBankCardForm(t *testing.T) {
	m := New(stubSaver{})
	// Press Enter on the first item (Bank Card)
	m = sendSpecial(m, tea.KeyEnter)
	view := m.View()
	if !strings.Contains(view, "Bank Card") && !strings.Contains(view, "💳") {
		t.Fatalf("expected Bank Card form after Enter; got:\n%s", view)
	}
}

func TestNavigateDownToCredentials_ThenEnter(t *testing.T) {
	m := New(stubSaver{})
	// Move to 2nd item (Credentials) and open
	m = sendSpecial(m, tea.KeyDown)
	m = sendSpecial(m, tea.KeyEnter)
	view := m.View()
	if !strings.Contains(view, "Credentials") && !strings.Contains(view, "🔑") {
		t.Fatalf("expected Credentials form; got:\n%s", view)
	}
}

func TestNavigateDownToTextData_ThenEnter(t *testing.T) {
	m := New(stubSaver{})
	// Move to 3rd item (Text Data) and open
	m = sendSpecial(m, tea.KeyDown)
	m = sendSpecial(m, tea.KeyDown)
	m = sendSpecial(m, tea.KeyEnter)
	view := m.View()
	if !strings.Contains(view, "Text Data") && !strings.Contains(view, "📝") {
		t.Fatalf("expected Text Data form; got:\n%s", view)
	}
}

func TestNavigateDownToBinaryData_ThenEnter(t *testing.T) {
	m := New(stubSaver{})
	// Move to 4th item (Binary Data) and open
	m = sendSpecial(m, tea.KeyDown)
	m = sendSpecial(m, tea.KeyDown)
	m = sendSpecial(m, tea.KeyDown)
	m = sendSpecial(m, tea.KeyEnter)
	view := m.View()
	if !strings.Contains(view, "Binary Data") && !strings.Contains(view, "📁") {
		t.Fatalf("expected Binary Data form; got:\n%s", view)
	}
}

// helper: send a special key (Up/Down/Enter/Space)
func sendSpecial(m Model, t tea.KeyType) Model {
	nm, _ := m.Update(tea.KeyMsg{Type: t})
	return nm.(Model)
}
