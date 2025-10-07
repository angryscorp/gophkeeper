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

// Sync orchestrates bidirectional synchronization between the local store
// (tokens/sync repository) and the remote gRPC service (client).
// It handles ping/auth check, pull (apply server changes), and push (flush outbox).
type Sync struct {
	client    Client
	tokenRepo Tokens
	syncRepo  Repository
}

// New creates a Sync use case with the provided remote client, token storage
// and local synchronization repository.
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

// Ping verifies the current session by calling the server's Ping method
// with the stored access token.
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

// Sync performs a full round:
//  1. Ping (auth check)
//  2. Pull (apply server changes)
//  3. Push (send local outbox)
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

// Pull fetches server changes page-by-page starting from the local cursor,
// applies them to the local store, and advances the cursor until exhausted.
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

// Push drains the local outbox in batches and sends them to the server.
// Successfully applied messages are removed from the outbox.
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
