package sync

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const prefixSyncService = "/gophkeeper.v1.SyncService/"

type TokenVerifier interface {
	Verify(token string) error
}

func UnaryAuthForSync(verifier TokenVerifier) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if !strings.HasPrefix(info.FullMethod, prefixSyncService) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("missing metadata")
		}

		h := md.Get("Authorization")
		if len(h) == 0 || !strings.HasPrefix(h[0], "Bearer ") {
			return nil, errors.New("missing bearer token")
		}

		token := strings.TrimSpace(h[0][len("Bearer "):])

		err := verifier.Verify(token)
		if err != nil {
			return nil, err
		}
		
		return handler(ctx, req)
	}
}
