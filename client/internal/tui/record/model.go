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
	infoType    domain.UserInfoType
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
}

func New() Model {
	return Model{
		route:  routeTypeSelection,
		cursor: 0,
		menuItems: []menuItem{
			{
				title:       "💳 Bank Card",
				description: "Credit/Debit card information",
				infoType:    domain.UserInfoTypeBankCard,
				route:       routeBankCardForm,
			},
			{
				title:       "🔑 Credentials",
				description: "Username and password",
				infoType:    domain.UserInfoTypeCredentials,
				route:       routeCredentialsForm,
			},
			{
				title:       "📝 Text Data",
				description: "Notes and text information",
				infoType:    domain.UserInfoTypeUserTextData,
				route:       routeTextDataForm,
			},
			{
				title:       "📁 Binary Data",
				description: "Files and binary information",
				infoType:    domain.UserInfoTypeUserBinaryData,
				route:       routeBinaryDataForm,
			},
		},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
