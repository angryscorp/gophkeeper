package register

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	signup   func(username, password string) error
	state    state
	input    textinput.Model
	err      error
	username string
	password string
}

type state int

const (
	stateInit state = iota
	stateAskUsername
	stateAskPassword
	stateInProgress
	stateSuccess
	stateError
)

func (m Model) Init() tea.Cmd {
	return cmdAskUsername.Run
}

func New(signup func(username, password string) error) Model {
	return Model{
		signup: signup,
		state:  stateInit,
		input:  textinput.New(),
	}
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
			return m, m.doRegister(m.username, m.password)
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
	title := "REGISTER"
	switch m.state {
	case stateAskUsername, stateAskPassword:
		return fmt.Sprintf("%s\n\n%s", title, m.input.View())
	case stateInProgress:
		return fmt.Sprintf("%s\n\nsending request...", title)
	case stateSuccess:
		return fmt.Sprintf("%s\n\nRegistration is successful", title)
	case stateError:
		return fmt.Sprintf("%s\n\nRegistration failed: %v", title, m.err)
	default:
		return fmt.Sprintf("%s\n\ninitializing...", title)
	}
}

func (m Model) doRegister(username, password string) tea.Cmd {
	return func() tea.Msg {
		err := m.signup(username, password)
		return resultMsg{
			success: err == nil,
			err:     err,
		}
	}
}
