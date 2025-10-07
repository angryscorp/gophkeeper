package sync

// state represents the current step of the sync model's lifecycle.
type state int

const (
	stateInit state = iota
	stateInProgress
	stateSuccess
	stateError
)
