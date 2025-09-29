package auth

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

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
				m.state = stateAskPassword
				return m, cmdAskPassword.Run
			}

			if m.state == stateAskPassword {
				m.password = m.input.Value()
				if m.password == "" {
					m.err = fmt.Errorf("password cannot be empty")
					m.state = stateError
					return m, nil
				}
				m.state = stateInProgress
				return m, cmdSendRequest.Run
			}

		default:
			if m.state == stateAskUsername || m.state == stateAskPassword {
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

		case cmdAskPassword:
			m.input.Prompt = "Enter password: "
			m.input.EchoMode = textinput.EchoPassword
			m.input.SetValue("")
			m.state = stateAskPassword
			return m, m.input.Focus()

		case cmdSendRequest:
			m.state = stateInProgress
			return m, m.doAction(m.username, m.password)
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

func (m Model) doAction(username, password string) tea.Cmd {
	return func() tea.Msg {
		err := m.action(username, password)
		return resultMsg{
			success: err == nil,
			err:     err,
		}
	}
}
