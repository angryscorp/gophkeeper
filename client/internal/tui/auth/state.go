package auth

type state int

const (
	stateInit state = iota
	stateAskUsername
	stateAskPassword
	stateInProgress
	stateSuccess
	stateError
)
