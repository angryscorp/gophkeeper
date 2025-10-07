package domain

// Error is a custom error type used across the domain layer.
type Error string

// Error implements the error interface.
func (e Error) Error() string {
	return string(e)
}

const (
	// ErrUsernameTaken is returned when trying to register an already existing username.
	ErrUsernameTaken Error = "username is taken"
	// ErrUsernameNotFound is returned when a username cannot be found.
	ErrUsernameNotFound Error = "username not found"
	// ErrChallengeNotFound is returned when a login challenge cannot be found (e.g., expired).
	ErrChallengeNotFound Error = "challenge not found"
	// ErrChallengeFailed is returned when a login challenge response is invalid.
	ErrChallengeFailed Error = "challenge failed"
)
