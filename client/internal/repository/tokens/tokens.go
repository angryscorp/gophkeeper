package tokens

import "context"

type Tokens interface {
	GetAccessToken(ctx context.Context) (string, error)
}
