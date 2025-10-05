package sync

import "context"

type Client interface {
	Ping(ctx context.Context, accessToken string) error
}
