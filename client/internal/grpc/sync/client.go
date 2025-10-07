package sync

import (
	"context"
	"errors"

	"gophkeeper/client/internal/domain"
	usecase "gophkeeper/client/internal/usecase/sync"
	"gophkeeper/pkg/grpc/sync"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

const pullBatchSize = 100

// Client is a wrapper around the gRPC SyncServiceClient.
// It provides higher-level methods for synchronization:
// Ping, Pull and Push, mapping gRPC messages to domain models.
type Client struct {
	client sync.SyncServiceClient
}

// New creates a new Sync client bound to the given gRPC connection.
func New(conn *grpc.ClientConn) *Client {
	return &Client{client: sync.NewSyncServiceClient(conn)}
}

var _ usecase.Client = (*Client)(nil)

// Ping calls the SyncService.Ping RPC to verify that the server
// is reachable and the provided access token is valid.
func (c Client) Ping(ctx context.Context, accessToken string) error {
	ctx = addTokenToContext(ctx, accessToken)
	_, err := c.client.Ping(ctx, &emptypb.Empty{})
	return err
}

// Pull calls the SyncService.Pull RPC to fetch batched changes
// from the server starting after the given cursor. It returns
// the list of messages, the updated cursor, and whether more
// changes are available.
func (c Client) Pull(ctx context.Context, accessToken string, cursor int64) (*usecase.PullResponse, error) {
	ctx = addTokenToContext(ctx, accessToken)
	req := &sync.PullRequest{Cursor: cursor, Limit: pullBatchSize}
	resp, err := c.client.Pull(ctx, req)
	if err != nil {
		return nil, err
	}

	res := make([]domain.Message, len(resp.Changes))
	for i, change := range resp.Changes {
		operationID, err := uuid.Parse(change.OperationId)
		if err != nil {
			return nil, err
		}

		recordID, err := uuid.Parse(change.Id)
		if err != nil {
			return nil, err
		}

		res[i] = domain.Message{
			ID:            operationID,
			RecordID:      recordID,
			Kind:          change.Kind,
			UpdatedAtUnix: change.UpdatedAtUnix,
			Payload:       change.Payload,
		}
	}

	return &usecase.PullResponse{Changes: res, NewCursor: resp.NextCursor, HasMore: resp.HasMore}, nil
}

// Push calls the SyncService.Push RPC to upload local changes
// to the server. It returns the list of record IDs that were
// successfully applied or acknowledged.
func (c Client) Push(ctx context.Context, accessToken string, messages []domain.Message) ([]uuid.UUID, error) {
	ctx = addTokenToContext(ctx, accessToken)

	changes := make([]*sync.RecordChange, len(messages))
	for i, message := range messages {
		changes[i] = &sync.RecordChange{
			OperationId:   message.ID.String(),
			Id:            message.RecordID.String(),
			Kind:          message.Kind,
			UpdatedAtUnix: message.UpdatedAtUnix,
			Payload:       message.Payload,
		}
	}

	resp, err := c.client.Push(ctx, &sync.PushRequest{Changes: changes})
	if err != nil {
		return nil, err
	}

	var res []uuid.UUID
	for _, result := range resp.Results {
		switch result.Status {
		case 1, 2:
			id, err := uuid.Parse(result.RecordId)
			if err != nil {
				return nil, err
			}
			res = append(res, id)
		default:
			return nil, errors.New("unknown result status")
		}
	}

	return res, err
}

// addTokenToContext attaches a bearer access token into
// the gRPC outgoing context metadata for authorization.
func addTokenToContext(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(
		ctx,
		"Authorization", "Bearer "+token,
	)
}
