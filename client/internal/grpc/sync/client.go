package sync

import (
	"context"
	"gophkeeper/pkg/grpc/sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	client sync.SyncServiceClient
}

func New(conn *grpc.ClientConn) *Client {
	return &Client{client: sync.NewSyncServiceClient(conn)}
}

func (c Client) Ping(ctx context.Context, accessToken string) error {
	ctx = addTokenToContext(ctx, accessToken)
	_, err := c.client.Ping(ctx, &emptypb.Empty{})
	return err
}

func addTokenToContext(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(
		ctx,
		"Authorization", "Bearer "+token,
	)
}
