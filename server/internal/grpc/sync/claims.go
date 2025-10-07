package sync

import "context"

type ctxKeyUsername struct{}

func usernameFromCtx(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyUsername{}).(string)
	return v, ok
}
