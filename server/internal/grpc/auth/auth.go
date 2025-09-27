package auth

import (
	"context"
	"gophkeeper/pkg/grpc/auth"
	"gophkeeper/pkg/grpc/mapper"
	usecaseAuth "gophkeeper/server/internal/usecase/auth"
)

type Server struct {
	auth.UnimplementedAuthServiceServer
	usecase *usecaseAuth.Auth
}

func New(usecase *usecaseAuth.Auth) *Server {
	return &Server{usecase: usecase}
}

var _ auth.AuthServiceServer = (*Server)(nil)

func (s *Server) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	err := s.usecase.Register(ctx, requestToDomain(req))
	if err != nil {
		return nil, err
	}
	return &auth.RegisterResponse{}, nil
}

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

func (s *Server) LoginFinish(ctx context.Context, req *auth.LoginFinishRequest) (*auth.LoginFinishResponse, error) {
	token, err := s.usecase.LoginFinish(ctx, req.Username, req.DeviceId, req.Response)
	if err != nil {
		return nil, err
	}
	return &auth.LoginFinishResponse{AccessToken: token}, nil
}
