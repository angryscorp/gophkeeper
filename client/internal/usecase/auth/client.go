package auth

import (
	"context"
	"gophkeeper/pkg/crypto"
)

type Client interface {
	Register(ctx context.Context, username string, kdf crypto.KDFParameters, edKey, authKey []byte, algorithm crypto.AuthKeyAlgorithm) error
	LoginStart(ctx context.Context, username string, deviceName string) (*crypto.LoginPayload, error)
	LoginFinish(ctx context.Context, username, deviceName string, challenge []byte) (string, error)
}
