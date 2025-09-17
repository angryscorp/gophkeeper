package menu

import (
	"fmt"
	"gophkeeper/client/internal/tui/auth"
	"gophkeeper/client/internal/tui/help"
	"gophkeeper/client/internal/tui/list"
	"gophkeeper/client/internal/tui/record"
	"gophkeeper/client/internal/tui/register"
	"gophkeeper/client/internal/tui/sync"
	usecaseAuth "gophkeeper/client/internal/usecase/auth"

	tea "github.com/charmbracelet/bubbletea"
)

type route int

const (
	routeMenu route = iota
	routeRegister
	routeAuth
	routeSync
	routeData
	routeNewItem
	routeHelp
	routeQuit
)

type menuItem struct {
	title string
	route route
	init  func(*Model) tea.Cmd
}

type Model struct {
	route  route
	items  []menuItem
	cursor int
	reg    register.Model
	auth   auth.Model
	sync   sync.Model
	data   list.Model
	record record.Model
	help   help.Model
}

func New(
	regFactory func() *usecaseAuth.Auth,
) Model {
	return Model{
		route: routeMenu,
		items: []menuItem{
			{"Register", routeRegister, func(m *Model) tea.Cmd { m.reg = register.New(regFactory()); return m.reg.Init() }},
			{"Login", routeAuth, func(m *Model) tea.Cmd { m.auth = auth.New(); return nil }},
			{"Sync", routeSync, func(m *Model) tea.Cmd { m.sync = sync.New(); return nil }},
			{"Private Data", routeData, func(m *Model) tea.Cmd { m.data = list.New(nil); return tea.WindowSize() }},
			{"Add New Item", routeNewItem, func(m *Model) tea.Cmd { m.record = record.New(); return nil }},
			{"Help", routeHelp, func(m *Model) tea.Cmd { m.help = help.New(); return nil }},
			{"Quit", routeQuit, func(m *Model) tea.Cmd { return tea.Quit }},
		},
		cursor: 1,
	}
}

func (m Model) view() string {
	str := "🔐 GophKeeper\n\n"
	for i, it := range m.items {
		cursor := " "
		if m.cursor == (i + 1) {
			cursor = "›"
		}
		str += fmt.Sprintf(" %s %s\n", cursor, it.title)
	}
	return str
}

func (m Model) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up":
			if m.cursor > 1 {
				m.cursor--
			}

		case "down":
			if m.cursor < len(m.items) {
				m.cursor++
			}

		case "enter":
			m.route = m.items[m.cursor-1].route
			cmd := m.items[m.cursor-1].init(&m)
			if cmd != nil {
				return m, cmd
			}
		}
	}

	return m, nil
}
