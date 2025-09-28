package record

import "fmt"

func (m Model) View() string {
	switch m.route {
	case routeTypeSelection:
		return m.view()
	case routeBankCardForm:
		return m.bankCard.View()
	case routeCredentialsForm:
		return m.credentials.View()
	case routeTextDataForm:
		return m.textData.View()
	case routeBinaryDataForm:
		return m.binaryData.View()
	default:
		return ""
	}
}

func (m Model) view() string {
	s := "🔐 Add New Item - Select Type\n\n"

	for i, item := range m.menuItems {
		cursor := "  "
		if m.cursor == i {
			cursor = "👉 "
		}
		s += fmt.Sprintf(" %s %s\n    %s\n\n", cursor, item.title, item.description)
	}

	s += "\n(use ↑/↓ to navigate, enter to select, ←/q/esc to return)"
	return s
}
