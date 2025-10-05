package auth

import "context"

type Tokens interface {
	Unlock(dataKey []byte) error
	SaveAccessToken(ctx context.Context, token string) error
}
