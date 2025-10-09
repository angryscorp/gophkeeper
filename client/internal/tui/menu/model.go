package menu

import (
	"gophkeeper/client/internal/tui/auth"
	"gophkeeper/client/internal/tui/help"
	"gophkeeper/client/internal/tui/list"
	"gophkeeper/client/internal/tui/record"
	"gophkeeper/client/internal/tui/sync"

	tea "github.com/charmbracelet/bubbletea"
)

type menuItem struct {
	title       string
	description string
	route       route
	init        func(*Model) tea.Cmd
}

// Model is the root Bubble Tea model that manages the main
// application menu. It holds sub-models for registration,
// login, sync, data listing, record editing, and help screens.
type Model struct {
	route  route
	items  []menuItem
	cursor int

	reg    auth.Model
	auth   auth.Model
	sync   sync.Model
	data   list.Model
	record record.Model
	help   help.Model
}

// New creates a new root menu Model with all menu items
// initialized. Each menu item is associated with a route
// and a factory function that constructs the corresponding
// sub-model.
func New(env Environment) Model {
	return Model{
		route: routeMenu,
		items: []menuItem{
			{
				"📕 Register",
				"Register new account",
				routeRegister,
				func(m *Model) tea.Cmd {
					m.reg = auth.New("📕 Register new account", env.RegFactory)
					return m.reg.Init()
				},
			},
			{
				"🔑️ Login",
				"Login to existing account",
				routeAuth,
				func(m *Model) tea.Cmd { m.auth = auth.New("🔑️ Login", env.LoginFactory); return m.auth.Init() },
			},
			{
				"📡 Sync",
				"Sync data with server",
				routeSync,
				func(m *Model) tea.Cmd { m.sync = sync.New(env.SyncFactory); return m.sync.Init() },
			},
			{
				"📒 Private Data",
				"See the vault's content",
				routeData,
				func(m *Model) tea.Cmd { m.data = list.New(env.DataFactory); return m.data.Init() },
			},
			{
				"💳 Add New Item",
				"Add data to the vault",
				routeNewItem,
				func(m *Model) tea.Cmd { m.record = record.New(env.DataSaver); return m.record.Init() },
			},
			{
				"📚 Help",
				"Help, controls, description",
				routeHelp,
				func(m *Model) tea.Cmd { m.help = help.New(env.HelpFactory); return m.help.Init() },
			},
			{
				"👋 Quit",
				"Exit Gophkeeper",
				routeQuit,
				func(m *Model) tea.Cmd { return tea.Quit },
			},
		},
		cursor: 1,
	}
}

// Init implements tea.Model. It does not perform any initialization
// because the menu is immediately ready to render.
func (m Model) Init() tea.Cmd {
	return nil
}
