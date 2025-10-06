package sync

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const prefixSyncService = "/gophkeeper.v1.SyncService/"

type TokenVerifier interface {
	Verify(token string) (string, error)
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
			return nil, status.Error(codes.Unauthenticated, "missing bearer token")
		}

		token := strings.TrimSpace(h[0][len("Bearer "):])

		username, err := verifier.Verify(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		ctx = context.WithValue(ctx, ctxKeyUsername{}, username)
		return handler(ctx, req)
	}
}
