package auth

import (
	"context"

	"gophkeeper/pkg/crypto"
)

type stubClient struct {
	registerFn    func(ctx context.Context, username string, kdf crypto.KDFParameters, edKey, authKey []byte, algo crypto.AuthKeyAlgorithm) error
	loginStartFn  func(ctx context.Context, username, deviceName string) (*crypto.LoginPayload, error)
	loginFinishFn func(ctx context.Context, username, deviceName string, challenge []byte) (string, error)

	lastRegister struct {
		username string
		kdf      crypto.KDFParameters
		edKey    []byte
		authKey  []byte
		algo     crypto.AuthKeyAlgorithm
	}
	lastLoginStart  struct{ username, deviceName string }
	lastLoginFinish struct {
		username, deviceName string
		challenge            []byte
	}
}

func (s *stubClient) Register(ctx context.Context, username string, kdf crypto.KDFParameters, edKey, authKey []byte, algo crypto.AuthKeyAlgorithm) error {
	s.lastRegister = struct {
		username string
		kdf      crypto.KDFParameters
		edKey    []byte
		authKey  []byte
		algo     crypto.AuthKeyAlgorithm
	}{username, kdf, edKey, authKey, algo}
	if s.registerFn != nil {
		return s.registerFn(ctx, username, kdf, edKey, authKey, algo)
	}
	return nil
}

func (s *stubClient) LoginStart(ctx context.Context, username string, deviceName string) (*crypto.LoginPayload, error) {
	s.lastLoginStart = struct{ username, deviceName string }{username, deviceName}
	if s.loginStartFn != nil {
		return s.loginStartFn(ctx, username, deviceName)
	}
	return &crypto.LoginPayload{}, nil
}

func (s *stubClient) LoginFinish(ctx context.Context, username, deviceName string, challenge []byte) (string, error) {
	s.lastLoginFinish = struct {
		username, deviceName string
		challenge            []byte
	}{username, deviceName, challenge}
	if s.loginFinishFn != nil {
		return s.loginFinishFn(ctx, username, deviceName, challenge)
	}
	return "TOKEN", nil
}
