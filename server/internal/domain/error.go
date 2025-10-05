package domain

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrUsernameTaken     Error = "username is taken"
	ErrUsernameNotFound  Error = "username not found"
	ErrChallengeNotFound Error = "challenge not found"
	ErrChallengeFailed   Error = "challenge failed"
)
