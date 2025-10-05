package list

type state int

const (
	stateInit state = iota
	stateInProgress
	stateSuccess
	stateError
)
