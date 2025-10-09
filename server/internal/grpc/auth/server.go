package auth

import (
	"context"

	"gophkeeper/pkg/grpc/auth"
	"gophkeeper/pkg/grpc/mapper"
	usecaseAuth "gophkeeper/server/internal/usecase/auth"
)

// Server implements the gRPC AuthServiceServer interface.
// It adapts the auth use case to gRPC transport.
type Server struct {
	auth.UnimplementedAuthServiceServer
	usecase *usecaseAuth.Auth
}

// New creates a new gRPC Auth server bound to the given use case.
func New(usecase *usecaseAuth.Auth) *Server {
	return &Server{usecase: usecase}
}

var _ auth.AuthServiceServer = (*Server)(nil)

// Register handles user registration requests.
func (s *Server) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	err := s.usecase.Register(ctx, requestToDomain(req))
	if err != nil {
		return nil, err
	}
	return &auth.RegisterResponse{}, nil
}

// LoginStart starts the login process by returning KDF parameters,
// encrypted data key, auth algorithm, and a challenge.
func (s *Server) LoginStart(ctx context.Context, req *auth.LoginStartRequest) (*auth.LoginStartResponse, error) {
	resp, err := s.usecase.LoginStart(ctx, req.Username, req.DeviceName)
	if err != nil {
		return nil, err
	}
	return &auth.LoginStartResponse{
		DeviceId:         resp.DeviceId,
		Kdf:              mapper.KdfParametersToGRPC(resp.KDFParameters),
		EncryptedDataKey: resp.EncryptedDataKey,
		AuthKeyAlg:       mapper.AuthAlgoToGRPC(resp.AuthKeyAlgorithm),
		Challenge:        resp.Challenge,
	}, nil
}

// LoginFinish completes the login process and returns an access token.
func (s *Server) LoginFinish(ctx context.Context, req *auth.LoginFinishRequest) (*auth.LoginFinishResponse, error) {
	token, err := s.usecase.LoginFinish(ctx, req.Username, req.DeviceId, req.Response)
	if err != nil {
		return nil, err
	}
	return &auth.LoginFinishResponse{AccessToken: token}, nil
}
