package auth

import "context"

// Tokens defines secure local storage for authentication data.
// It is responsible for unlocking the encrypted database with
// a data key and persisting issued access tokens.
type Tokens interface {
	Unlock(dataKey []byte) error
	SaveAccessToken(ctx context.Context, token string) error
}
