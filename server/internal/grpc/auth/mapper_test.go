package auth

import (
	"reflect"
	"testing"

	"gophkeeper/pkg/crypto"
	"gophkeeper/pkg/grpc/auth"
	"gophkeeper/pkg/grpc/mapper"

	"github.com/google/uuid"
)

func TestRequestToDomain(t *testing.T) {
	req := &auth.RegisterRequest{
		Username:         "alice",
		Kdf:              mapper.KdfParametersToGRPC(crypto.DefaultKDFParameters()),
		EncryptedDataKey: []byte{1, 2, 3},
		AuthKeyAlg:       auth.AuthKeyAlg_HMAC_SHA256,
		AuthKey:          []byte{9, 9, 9},
	}

	user := requestToDomain(req)

	if user.Username != req.Username {
		t.Errorf("expected username %q, got %q", req.Username, user.Username)
	}
	if !reflect.DeepEqual(user.EncryptedDataKey, req.EncryptedDataKey) {
		t.Errorf("expected EncryptedDataKey %v, got %v", req.EncryptedDataKey, user.EncryptedDataKey)
	}
	if !reflect.DeepEqual(user.AuthKey, req.AuthKey) {
		t.Errorf("expected AuthKey %v, got %v", req.AuthKey, user.AuthKey)
	}
	if string(user.AuthKeyAlgorithm) != "HMAC_SHA256" {
		t.Errorf("expected AuthKeyAlgorithm HMAC_SHA256, got %s", user.AuthKeyAlgorithm)
	}
	if user.ID == uuid.Nil {
		t.Errorf("expected non-nil UUID, got Nil")
	}
	if user.KDFParameters.TimeCost != req.Kdf.TimeCost {
		t.Errorf("expected TimeCost %d, got %d", req.Kdf.TimeCost, user.KDFParameters.TimeCost)
	}
}
