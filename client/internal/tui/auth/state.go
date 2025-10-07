package auth

// state represents the current step in the authentication flow.
// It drives the TUI state machine, from initial screen through
// input collection to request execution and result handling.
type state int

const (
	stateInit state = iota
	stateAskUsername
	stateAskPassword
	stateInProgress
	stateSuccess
	stateError
)
