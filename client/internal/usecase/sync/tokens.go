package sync

import "context"

// Tokens provides access to locally stored authentication tokens.
type Tokens interface {
	GetAccessToken(ctx context.Context) (string, error)
	Ready() bool
}
