package list

import (
	"gophkeeper/client/internal/domain"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	list        list.Model
	dataFactory func() ([]domain.Record, error)
	state       state
	err         error
}

func New(dataFactory func() ([]domain.Record, error)) Model {
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = "📒 Private Data"
	return Model{list: l, dataFactory: dataFactory}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.WindowSize(),
		cmdLoadData.Run,
	)
}
