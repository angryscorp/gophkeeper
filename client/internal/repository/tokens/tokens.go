package tokens

import "context"

const (
	AccessTokenNotFound = "access token not found"
)

type Tokens interface {
	Unlock(dataKey []byte) error
	GetAccessToken(ctx context.Context) (string, error)
	SaveAccessToken(ctx context.Context, token string) error
}
