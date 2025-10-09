package common

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// InputWithPlaceholder returns a textinput.Model preconfigured with the given
// placeholder and initial width matching the placeholder length.
func InputWithPlaceholder(placeholder string) textinput.Model {
	input := textinput.New()
	input.Placeholder = placeholder
	input.Width = len(placeholder)
	return input
}

// ExtendInputWidthIfNeeded grows the input's width to fit the current value
// length, if the value is wider than the configured width.
func ExtendInputWidthIfNeeded(input *textinput.Model) {
	currentWidth := len(input.Value())
	if input.Width < currentWidth {
		input.Width = currentWidth
	}
}

// UpdateSubModel calls Update(msg) on a Bubble Tea submodel if it implements
// the expected interface, returning the updated submodel (typed as T) and the
// resulting tea.Cmd. If the submodel does not implement Update, it returns the
// original submodel and a nil command.
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
