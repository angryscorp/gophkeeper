package list

import (
	"gophkeeper/client/internal/domain"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// Model is a Bubble Tea model for browsing and displaying
// a list of stored records (credentials, cards, notes, etc.)
// retrieved via the provided data factory function.
type Model struct {
	list        list.Model
	dataFactory func() ([]domain.Record, error)
	state       state
	err         error
}

// New creates a new list Model with a title and the given
// data factory callback, which is used to load records.
func New(dataFactory func() ([]domain.Record, error)) Model {
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = "📒 Private Data"
	return Model{list: l, dataFactory: dataFactory}
}

// Init implements tea.Model. It initializes the list by
// requesting the window size and loading data.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.WindowSize(),
		cmdLoadData.Run,
	)
}
