package sync

import (
	"context"
	"fmt"
	"gophkeeper/pkg/grpc/sync"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	sync.UnimplementedSyncServiceServer
}

func New() *Server {
	return &Server{}
}

var _ sync.SyncServiceServer = (*Server)(nil)

func (s Server) Pull(ctx context.Context, request *sync.PullRequest) (*sync.PullResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) Push(ctx context.Context, request *sync.PushRequest) (*sync.PushResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) Ping(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	fmt.Println("ping")
	return &emptypb.Empty{}, nil
}
