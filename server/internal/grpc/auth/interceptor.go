package auth

import (
	"context"
	"errors"
	"gophkeeper/server/internal/domain"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorMappingServerInterceptor is a gRPC interceptor that maps
// domain-level errors into appropriate gRPC status codes.
func ErrorMappingServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		resp, err = handler(ctx, req)
		return resp, mapError(err)
	}
}

func mapError(err error) error {
	if err == nil {
		return nil
	}

	errorCode := codes.Unknown
	switch {
	case errors.Is(err, context.Canceled):
		errorCode = codes.Canceled
	case errors.Is(err, domain.ErrUsernameTaken):
		errorCode = codes.AlreadyExists
	case errors.Is(err, domain.ErrUsernameNotFound):
		errorCode = codes.NotFound
	case errors.Is(err, domain.ErrChallengeNotFound):
		errorCode = codes.FailedPrecondition
	case errors.Is(err, domain.ErrChallengeFailed):
		errorCode = codes.Unauthenticated
	}

	return status.Error(errorCode, err.Error())
}
