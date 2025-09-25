package tokens

import "context"

type Tokens interface {
	GetAccessToken(ctx context.Context) (string, error)
	SaveAccessToken(ctx context.Context, token string) error
}
