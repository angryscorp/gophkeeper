package record

import (
	"gophkeeper/client/internal/domain"
	"gophkeeper/client/internal/tui/bankcard"
	"gophkeeper/client/internal/tui/binarydata"
	"gophkeeper/client/internal/tui/credentials"
	"gophkeeper/client/internal/tui/textdata"

	tea "github.com/charmbracelet/bubbletea"
)

type menuItem struct {
	title       string
	description string
	route       recordRoute
}

type Model struct {
	route     recordRoute
	cursor    int
	menuItems []menuItem

	bankCard    bankcard.Model
	credentials credentials.Model
	textData    textdata.Model
	binaryData  binarydata.Model

	dataSaver domain.UserInfoSaver
}

func New(dataSaver domain.UserInfoSaver) Model {
	return Model{
		route:  routeTypeSelection,
		cursor: 0,
		menuItems: []menuItem{
			{
				title:       domain.UserDataKindBankCard.Title(),
				description: domain.UserDataKindBankCard.Description(),
				route:       routeBankCardForm,
			},
			{
				title:       domain.UserDataKindCredentials.Title(),
				description: domain.UserDataKindCredentials.Description(),
				route:       routeCredentialsForm,
			},
			{
				title:       domain.UserDataKindTextData.Title(),
				description: domain.UserDataKindTextData.Description(),
				route:       routeTextDataForm,
			},
			{
				title:       domain.UserDataKindBinaryData.Title(),
				description: domain.UserDataKindBinaryData.Description(),
				route:       routeBinaryDataForm,
			},
		},
		dataSaver: dataSaver,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
