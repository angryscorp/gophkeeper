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
	client     Client
	tokenRepo  Tokens
	outboxRepo OutboxRepository
}

func New(
	client Client,
	tokenRepo Tokens,
	outboxRepo OutboxRepository,
) *Sync {
	return &Sync{
		client:     client,
		tokenRepo:  tokenRepo,
		outboxRepo: outboxRepo,
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

func (sync *Sync) Push() error {
	if !sync.tokenRepo.Ready() {
		return errors.New("no active login session")
	}

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	for {
		messages, err := sync.outboxRepo.GetBatch(ctx, batchLimit)
		if err != nil {
			return err
		}

		if len(messages) == 0 {
			return nil
		}

		token, err := sync.tokenRepo.GetAccessToken(ctx)
		if err != nil {
			return err
		}

		resp, err := sync.client.Push(ctx, token, messages)
		if err != nil {
			return err
		}

		err = sync.outboxRepo.DeleteBatch(ctx, resp)
		if err != nil {
			return err
		}
	}
}
