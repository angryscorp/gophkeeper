package tokens

import "context"

type Tokens interface {
	Unlock(dataKey []byte) error
	GetAccessToken(ctx context.Context) (string, error)
	SaveAccessToken(ctx context.Context, token string) error
	Ready() bool
}
