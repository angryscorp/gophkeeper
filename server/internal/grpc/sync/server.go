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
	fmt.Printf("pull: %v\n", request)
	return &sync.PullResponse{}, nil
}

func (s Server) Push(ctx context.Context, request *sync.PushRequest) (*sync.PushResponse, error) {
	fmt.Println("push")
	for i, rec := range request.Changes {
		fmt.Printf("->rec %d; id: %s; kind: %d\n", i, rec.Id, rec.Kind)
	}

	// TODO: saving logic

	if len(request.Changes) == 0 {
		return &sync.PushResponse{
			Results: []*sync.PushResult{},
		}, nil
	} else {
		return &sync.PushResponse{
			Results: []*sync.PushResult{
				{
					RecordId: request.Changes[0].OperationId,
					Status:   1,
				},
			},
		}, nil
	}
}

func (s Server) Ping(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	fmt.Println("ping")
	return &emptypb.Empty{}, nil
}
