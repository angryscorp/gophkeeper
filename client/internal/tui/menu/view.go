package menu

import "fmt"

// View implements tea.Model for the root menu. It renders either
// the main menu or the currently active sub-model screen depending
// on the current route.
func (m Model) View() string {
	switch m.route {
	case routeMenu:
		return m.view()
	case routeRegister:
		return m.reg.View()
	case routeAuth:
		return m.auth.View()
	case routeSync:
		return m.sync.View()
	case routeData:
		return m.data.View()
	case routeNewItem:
		return m.record.View()
	case routeHelp:
		return m.help.View()
	default:
		return ""
	}
}

func (m Model) view() string {
	str := "🔐 GophKeeper \n\n"
	for i, item := range m.items {
		cursor := "  "
		if m.cursor == (i + 1) {
			cursor = "👉 "
		}
		str += fmt.Sprintf(" %s %s\n    %s\n\n", cursor, item.title, item.description)
	}

	str += "\n(use ↑/↓ to navigate, enter to select, ←/q/esc to return)"
	return str
}
