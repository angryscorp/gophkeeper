package menu

import (
	"gophkeeper/client/internal/tui/common"

	tea "github.com/charmbracelet/bubbletea"
)

// Update implements tea.Model for the root menu. It routes incoming
// messages either to the main menu handler or to the currently active
// sub-model (register/login/sync/data/new item/help), and handles
// returning back to the menu on Esc/←/q.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.route {
	case routeMenu:
		return m.update(msg)

	case routeRegister, routeAuth, routeSync, routeData, routeNewItem, routeHelp:
		return handleSubModelUpdate(m, msg)

	default:
		return m, nil
	}
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

func handleSubModelUpdate(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	if isExitKey(msg) {
		m.route = routeMenu
		return m, nil
	}

	var cmd tea.Cmd
	switch m.route {
	case routeRegister:
		m.reg, cmd = common.UpdateSubModel(m.reg, msg)
	case routeAuth:
		m.auth, cmd = common.UpdateSubModel(m.auth, msg)
	case routeSync:
		m.sync, cmd = common.UpdateSubModel(m.sync, msg)
	case routeData:
		m.data, cmd = common.UpdateSubModel(m.data, msg)
	case routeNewItem:
		m.record, cmd = common.UpdateSubModel(m.record, msg)
	case routeHelp:
		m.help, cmd = common.UpdateSubModel(m.help, msg)
	default:
	}

	return m, cmd
}

func isExitKey(msg tea.Msg) bool {
	switch msg.(type) {
	case tea.KeyMsg:
		km := msg.(tea.KeyMsg)
		return km.Type == tea.KeyEsc || km.String() == "left" || km.String() == "q"
	}
	return false
}
