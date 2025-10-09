package sync

import "context"

type ctxKeyUsername struct{}

func usernameFromCtx(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyUsername{}).(string)
	return v, ok
}

func contextWithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, ctxKeyUsername{}, username)
}
