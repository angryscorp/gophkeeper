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

type Client struct {
	client sync.SyncServiceClient
}

func New(conn *grpc.ClientConn) *Client {
	return &Client{client: sync.NewSyncServiceClient(conn)}
}

var _ usecase.Client = (*Client)(nil)

func (c Client) Ping(ctx context.Context, accessToken string) error {
	ctx = addTokenToContext(ctx, accessToken)
	_, err := c.client.Ping(ctx, &emptypb.Empty{})
	return err
}

func (c Client) Pull(ctx context.Context, accessToken string, cursor int64) (*usecase.PullResponse, error) {
	ctx = addTokenToContext(ctx, accessToken)
	req := &sync.PullRequest{Cursor: cursor, Limit: pullBatchSize}
	resp, err := c.client.Pull(ctx, req)
	if err != nil {
		return nil, err
	}

	res := make([]domain.Message, len(resp.Changes))
	for i, change := range resp.Changes {
		operationID, err := uuid.Parse(change.Change.OperationId)
		if err != nil {
			return nil, err
		}

		recordID, err := uuid.Parse(change.Change.Id)
		if err != nil {
			return nil, err
		}

		res[i] = domain.Message{
			ID:            operationID,
			RecordID:      recordID,
			Kind:          change.Change.Kind,
			UpdatedAtUnix: change.Change.UpdatedAtUnix,
			Payload:       change.Change.Payload,
		}
	}

	return &usecase.PullResponse{Changes: res, NewCursor: resp.NextCursor, HasMore: resp.HasMore}, nil
}

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

func addTokenToContext(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(
		ctx,
		"Authorization", "Bearer "+token,
	)
}
