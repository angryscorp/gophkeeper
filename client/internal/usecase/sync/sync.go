package sync

import (
	"context"
	"errors"
	"time"
)

const (
	ctxTimeout = 5 * time.Second
	batchLimit = 100
)

type Sync struct {
	client    Client
	tokenRepo Tokens
	syncRepo  Repository
}

func New(
	client Client,
	tokenRepo Tokens,
	syncRepo Repository,
) *Sync {
	return &Sync{
		client:    client,
		tokenRepo: tokenRepo,
		syncRepo:  syncRepo,
	}
}

func (sync *Sync) Ping() error {
	if !sync.tokenRepo.Ready() {
		return errors.New("no active login session")
	}

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	token, err := sync.tokenRepo.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	return sync.client.Ping(ctx, token)
}

func (sync *Sync) Sync() error {
	err := sync.Ping()
	if err != nil {
		return err
	}

	err = sync.Pull()
	if err != nil {
		return err
	}

	err = sync.Push()
	if err != nil {
		return err
	}

	return nil
}

func (sync *Sync) Pull() error {
	if !sync.tokenRepo.Ready() {
		return errors.New("no active login session")
	}

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	token, err := sync.tokenRepo.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	for {
		cursor, err := sync.syncRepo.GetCursor(ctx)
		if err != nil {
			return err
		}

		resp, err := sync.client.Pull(ctx, token, cursor)
		if err != nil {
			return err
		}

		err = sync.syncRepo.SaveChanges(ctx, resp.Changes, resp.NewCursor)
		if err != nil {
			return err
		}

		if !resp.HasMore {
			break
		}
	}

	return nil
}

func (sync *Sync) Push() error {
	if !sync.tokenRepo.Ready() {
		return errors.New("no active login session")
	}

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	token, err := sync.tokenRepo.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	for {
		messages, err := sync.syncRepo.GetBatch(ctx, batchLimit)
		if err != nil {
			return err
		}

		if len(messages) == 0 {
			return nil
		}

		resp, err := sync.client.Push(ctx, token, messages)
		if err != nil {
			return err
		}

		err = sync.syncRepo.DeleteBatch(ctx, resp)
		if err != nil {
			return err
		}
	}
}
