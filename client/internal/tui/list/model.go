package list

import (
	"github.com/charmbracelet/bubbles/list"
)

type Model struct {
	List list.Model
}

type Item struct{ title, desc string }

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }

func New(items []list.Item) Model {
	if items == nil {
		// demo data
		items = []list.Item{
			Item{title: "Google Account", desc: "Credentials"},
			Item{title: "My Notes", desc: "Text data"},
			Item{title: "Swedbank", desc: "Bank Card"},
			Item{title: "Revolut", desc: "Bank Card"},
			Item{title: "My photo", desc: "Binary Data"},
		}
	}

	m := Model{List: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.List.SetShowTitle(false)

	return m
}
