package menu

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.route {
	case routeMenu:
		return m.update(msg)

	case routeRegister, routeAuth, routeSync, routeData, routeNewItem, routeHelp:
		return handleSubModelUpdate(&m, msg)

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

func handleSubModelUpdate(m *Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	if isExitKey(msg) {
		m.route = routeMenu
		return m, nil
	}

	var cmd tea.Cmd
	switch m.route {
	case routeRegister:
		m.reg, cmd = updateSubModel(m.reg, msg)
	case routeAuth:
		m.auth, cmd = updateSubModel(m.auth, msg)
	case routeSync:
		m.sync, cmd = updateSubModel(m.sync, msg)
	case routeData:
		m.data, cmd = updateSubModel(m.data, msg)
	case routeNewItem:
		m.record, cmd = updateSubModel(m.record, msg)
	case routeHelp:
		m.help, cmd = updateSubModel(m.help, msg)
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

func updateSubModel[T any](subModel T, msg tea.Msg) (T, tea.Cmd) {
	updater, ok := any(subModel).(interface {
		Update(tea.Msg) (tea.Model, tea.Cmd)
	})
	if !ok {
		return subModel, nil
	}

	newSub, cmd := updater.Update(msg)
	if nm, ok := newSub.(T); ok {
		return nm, cmd
	}
	return subModel, cmd
}
