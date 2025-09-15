package auth

import (
	"context"
	"time"

	"gophkeeper/pkg/grpc/auth"
)

type Server struct {
	auth.UnimplementedAuthServiceServer
}

func New() *Server {
	return &Server{}
}

var _ auth.AuthServiceServer = (*Server)(nil)

func (s *Server) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	return &auth.RegisterResponse{}, nil
}

func (s *Server) LoginStart(ctx context.Context, req *auth.LoginStartRequest) (*auth.LoginStartResponse, error) {
	return &auth.LoginStartResponse{
		DeviceId:         "dev-stub",
		Kdf:              nil,
		EncryptedDataKey: nil,
		AuthKeyAlg:       auth.AuthKeyAlg_HMAC_SHA256,
		Challenge:        []byte("1234"),
	}, nil
}

func (s *Server) LoginFinish(ctx context.Context, req *auth.LoginFinishRequest) (*auth.LoginFinishResponse, error) {
	return &auth.LoginFinishResponse{
		AccessToken:   "access-token",
		RefreshToken:  "refresh-token",
		ExpiresAtUnix: time.Now().Add(15 * time.Minute).Unix(),
	}, nil
}

func (s *Server) RefreshToken(ctx context.Context, req *auth.RefreshRequest) (*auth.RefreshResponse, error) {
	return &auth.RefreshResponse{
		AccessToken:   "access-token",
		ExpiresAtUnix: time.Now().Add(15 * time.Minute).Unix(),
	}, nil
}

func (s *Server) ChangePassword(ctx context.Context, req *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	return &auth.ChangePasswordResponse{}, nil
}
