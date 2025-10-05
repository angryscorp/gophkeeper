package domain

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrUsernameTaken    Error = "username already taken 🥲"
	ErrUsernameNotFound Error = "user not found 🤷‍♂️"
)
