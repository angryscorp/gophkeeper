package common

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func InputWithPlaceholder(placeholder string) textinput.Model {
	input := textinput.New()
	input.Placeholder = placeholder
	input.Width = len(placeholder)
	return input
}

func ExtendInputWidthIfNeeded(input *textinput.Model) {
	currentWidth := len(input.Value())
	if input.Width < currentWidth {
		input.Width = currentWidth
	}
}

func UpdateSubModel[T any](subModel T, msg tea.Msg) (T, tea.Cmd) {
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
