package record

import (
	"gophkeeper/client/internal/tui/bankcard"
	"gophkeeper/client/internal/tui/binarydata"
	"gophkeeper/client/internal/tui/common"
	"gophkeeper/client/internal/tui/credentials"
	"gophkeeper/client/internal/tui/textdata"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.route {
	case routeTypeSelection:
		return m.update(msg)
	case routeBankCardForm, routeCredentialsForm, routeTextDataForm, routeBinaryDataForm:
		return handleSubModelUpdate(m, msg)
	default:
		return m, nil
	}
}

func handleSubModelUpdate(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.route {
	case routeBankCardForm:
		m.bankCard, cmd = common.UpdateSubModel(m.bankCard, msg)
	case routeCredentialsForm:
		m.credentials, cmd = common.UpdateSubModel(m.credentials, msg)
	case routeTextDataForm:
		m.textData, cmd = common.UpdateSubModel(m.textData, msg)
	case routeBinaryDataForm:
		m.binaryData, cmd = common.UpdateSubModel(m.binaryData, msg)
	default:
	}

	return m, cmd
}
func (m Model) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.menuItems)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.route = m.menuItems[m.cursor].route
			return m.initSubModel()
		}
	}
	return m, nil
}

func (m Model) initSubModel() (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.route {
	case routeBankCardForm:
		m.bankCard = bankcard.New()
		cmd = m.bankCard.Init()
	case routeCredentialsForm:
		m.credentials = credentials.New(m.dataSaver.SaveCredentials)
		cmd = m.credentials.Init()
	case routeTextDataForm:
		m.textData = textdata.New()
		cmd = m.textData.Init()
	case routeBinaryDataForm:
		m.binaryData = binarydata.New()
		cmd = m.binaryData.Init()
	default:
	}

	return m, cmd
}
