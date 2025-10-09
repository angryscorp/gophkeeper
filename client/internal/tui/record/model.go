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

// Model is a Bubble Tea model for creating new records.
// It starts with a type-selection menu and routes into
// one of the data entry forms (bank card, credentials,
// text data, binary data).
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

// New creates a new record creation Model, wiring it
// with the provided data saver and initializing the
// menu of available record types.
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

// Init implements tea.Model. It performs no initialization.
func (m Model) Init() tea.Cmd {
	return nil
}
