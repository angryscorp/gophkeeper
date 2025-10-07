package sync

import (
	"context"

	"gophkeeper/pkg/grpc/sync"
	"gophkeeper/server/internal/domain"
	usecaseSync "gophkeeper/server/internal/usecase/sync"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Server implements the gRPC SyncServiceServer and
// adapts the sync use case to gRPC transport.
type Server struct {
	sync.UnimplementedSyncServiceServer
	usecase *usecaseSync.Sync
}

// New creates a new Sync gRPC server bound to the given use case.
func New(usecase *usecaseSync.Sync) *Server {
	return &Server{usecase: usecase}
}

var _ sync.SyncServiceServer = (*Server)(nil)

// Pull returns a batch of changes for the given username starting from the cursor.
func (s Server) Pull(ctx context.Context, request *sync.PullRequest) (*sync.PullResponse, error) {
	username, ok := usernameFromCtx(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no username in context")
	}

	resp, err := s.usecase.Pull(ctx, username, request.Cursor, request.Limit)
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

// Push stores a batch of changes from the client and returns per-record results.
func (s Server) Push(ctx context.Context, request *sync.PushRequest) (*sync.PushResponse, error) {
	username, ok := usernameFromCtx(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no username in context")
	}

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

	resp, err := s.usecase.Push(ctx, username, records)
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

// Ping is a simple liveness check.
func (s Server) Ping(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
