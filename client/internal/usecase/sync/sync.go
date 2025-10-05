package sync

import (
	"context"
	"errors"
	"time"
)

const (
	ctxTimeout = 5 * time.Second
)

type Sync struct {
	client Client
	repo   Tokens
}

func New(
	client Client,
	repo Tokens,
) *Sync {
	return &Sync{
		client: client,
		repo:   repo,
	}
}

func (sync *Sync) Ping() error {
	if !sync.repo.Ready() {
		return errors.New("no active login session")
	}

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	token, err := sync.repo.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	return sync.client.Ping(ctx, token)
}
