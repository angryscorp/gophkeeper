package register

import (
	"fmt"
	"gophkeeper/client/internal/usecase/auth"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	auth  *auth.Auth
	state state
	input textinput.Model
	err   error
	name  string
}

type state int

const (
	stateInit state = iota
	stateAskUsername
	stateInProgress
	stateSuccess
	stateError
)

func (m Model) Init() tea.Cmd {
	return cmdAskUsername.Run
}

func New(auth *auth.Auth) Model {
	return Model{
		auth:  auth,
		state: stateInit,
		input: func() textinput.Model {
			ti := textinput.New()
			ti.Prompt = "Enter username: "
			ti.CharLimit = 64
			return ti
		}(),
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.state == stateAskUsername {
				m.name = m.input.Value()
				if m.name == "" {
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
			m.state = stateAskUsername
			return m, m.input.Focus()

		case cmdSendRequest:
			m.state = stateInProgress
			return m, m.doRegister(m.name)
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
	switch m.state {
	case stateAskUsername:
		return fmt.Sprintf("REGISTER\n\n%s\n\n[Enter] continue • [q] quit", m.input.View())
	case stateInProgress:
		return "REGISTER\n\nsending request...\n\n[q] quit"
	case stateSuccess:
		return fmt.Sprintf("REGISTER\n\nok! username=%q\n\n[q] quit", m.name)
	case stateError:
		return fmt.Sprintf("REGISTER\n\nerror: %v\n\n[q] quit", m.err)
	default:
		return "REGISTER\n\n(initializing...)\n"
	}
}

func (m Model) doRegister(username string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second) // testing load state
		err := m.auth.Register(username)
		return resultMsg{
			success: err == nil,
			err:     err,
		}
	}
}
