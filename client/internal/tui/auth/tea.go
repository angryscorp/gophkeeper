package auth

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	login    func(username string) error
	state    state
	input    textinput.Model
	err      error
	username string
}

type state int

const (
	stateInit state = iota
	stateAskUsername
	stateInProgress
	stateSuccess
	stateError
)

func New(login func(username string) error) Model {
	return Model{
		login: login,
		state: stateInit,
		input: textinput.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return cmdAskUsername.Run
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.state == stateAskUsername {
				m.username = m.input.Value()
				if m.username == "" {
					m.err = fmt.Errorf("username cannot be empty")
					m.state = stateError
					return m, nil
				}
				m.state = stateInProgress
				return m, cmdSendRequest.Run
			}

		default:
			if m.state == stateAskUsername {
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)
				return m, cmd
			}
		}

	case baseCmd:
		switch msg {
		case cmdAskUsername:
			m.input.Prompt = "Enter username: "
			m.input.EchoMode = textinput.EchoNormal
			m.input.SetValue("")
			m.state = stateAskUsername
			return m, m.input.Focus()

		case cmdSendRequest:
			m.state = stateInProgress
			return m, m.doLogin(m.username)
		}

	case resultMsg:
		if msg.success {
			m.state = stateSuccess
		} else {
			m.err = msg.err
			m.state = stateError
		}
		return m, nil
	}

	return m, nil
}

func (m Model) View() string {
	title := "LOGIN"
	switch m.state {
	case stateAskUsername:
		return fmt.Sprintf("%s\n\n%s", title, m.input.View())
	case stateInProgress:
		return fmt.Sprintf("%s\n\nsending request...", title)
	case stateSuccess:
		return fmt.Sprintf("%s\n\nLogin is successful", title)
	case stateError:
		return fmt.Sprintf("%s\n\nLogin failed: %v", title, m.err)
	default:
		return fmt.Sprintf("%s\n\ninitializing...", title)
	}
}

func (m Model) doLogin(username string) tea.Cmd {
	return func() tea.Msg {
		err := m.login(username)
		return resultMsg{
			success: err == nil,
			err:     err,
		}
	}
}
