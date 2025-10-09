package domain

// Error is a custom error type used for domain-level errors.
// It implements the standard error interface.
type Error string

// Error returns the string message of the domain error.
func (e Error) Error() string {
	return string(e)
}

const (
	// ErrUsernameTaken indicates that the chosen username is already registered.
	ErrUsernameTaken Error = "username already taken 🥲"
	// ErrUsernameNotFound indicates that the specified username does not exist in the system.
	ErrUsernameNotFound Error = "user not found 🤷‍♂️"
)
