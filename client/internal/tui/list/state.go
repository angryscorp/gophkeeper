package list

// state represents the internal state of the list model.
// It is used to track progress of data loading and error handling.
type state int

const (
	stateInit state = iota
	stateInProgress
	stateSuccess
	stateError
)
