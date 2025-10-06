package sync

import (
	"context"
	"gophkeeper/pkg/grpc/sync"
	"gophkeeper/server/internal/domain"
	usecaseSync "gophkeeper/server/internal/usecase/sync"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	sync.UnimplementedSyncServiceServer
	usecase *usecaseSync.Sync
}

func New(usecase *usecaseSync.Sync) *Server {
	return &Server{usecase: usecase}
}

var _ sync.SyncServiceServer = (*Server)(nil)

func (s Server) Pull(ctx context.Context, request *sync.PullRequest) (*sync.PullResponse, error) {
	resp, err := s.usecase.Pull(ctx, request.Cursor, request.Limit)
	if err != nil {
		return nil, err
	}

	res := make([]*sync.RecordChange, len(resp.Changes))
	for i, change := range resp.Changes {
		res[i] = &sync.RecordChange{
			Id:            change.RecordID.String(),
			Kind:          change.Kind,
			UpdatedAtUnix: change.UpdatedAtUnix,
			Payload:       change.Payload,
			OperationId:   change.ID.String(),
		}
	}

	return &sync.PullResponse{
		Changes:    res,
		NextCursor: resp.NextCursor,
		HasMore:    resp.HasMore,
	}, nil
}

func (s Server) Push(ctx context.Context, request *sync.PushRequest) (*sync.PushResponse, error) {
	records := make([]domain.Message, len(request.Changes))
	for i, change := range request.Changes {
		operationID, err := uuid.Parse(change.OperationId)
		if err != nil {
			return nil, err
		}

		recordID, err := uuid.Parse(change.Id)
		if err != nil {
			return nil, err
		}

		records[i] = domain.Message{
			ID:            operationID,
			RecordID:      recordID,
			Kind:          change.Kind,
			UpdatedAtUnix: change.UpdatedAtUnix,
			Payload:       change.Payload,
		}
	}

	resp, err := s.usecase.Push(ctx, records)
	if err != nil {
		return nil, err
	}

	res := make([]*sync.PushResult, len(resp))
	for i, record := range resp {
		res[i] = &sync.PushResult{
			RecordId: record.String(),
			Status:   1,
		}
	}

	return &sync.PushResponse{Results: res}, nil
}

func (s Server) Ping(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
