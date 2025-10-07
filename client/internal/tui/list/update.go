package list

import (
	"gophkeeper/client/internal/domain"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i := m.list.Cursor()
			item := m.list.Items()[i].(item)
			item.showInfo = !item.showInfo
			m.list.Items()[i] = item
		}

	case baseCmd:
		switch msg {
		case cmdLoadData:
			m.state = stateInProgress
			return m, func() tea.Msg {
				data, err := m.dataFactory()
				return resultMsg{
					data: mapRecords(data),
					err:  err,
				}
			}
		}

	case resultMsg:
		if msg.err != nil {
			m.state = stateError
			m.err = msg.err
			return m, nil
		}
		m.state = stateSuccess
		return m, m.list.SetItems(msg.data)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func mapRecords(records []domain.Record) []list.Item {
	if records == nil || len(records) == 0 {
		return []list.Item{}
	}
	res := make([]list.Item, len(records))
	for i, r := range records {
		res[i] = item{
			title: r.Title,
			kind:  r.Kind.Title(),
			info:  r.SensitiveInfo,
		}
	}
	return res
}
