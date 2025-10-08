package auth

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"gophkeeper/pkg/crypto"
)

func TestRegister_Success(t *testing.T) {
	client := &stubClient{}
	tokens := &stubTokens{}
	a := New(client, tokens, func([]byte) {})

	if err := a.Register("alice", "password123"); err != nil {
		t.Fatalf("Register unexpected error: %v", err)
	}
	if client.lastRegister.username != "alice" {
		t.Errorf("username = %q, want alice", client.lastRegister.username)
	}
	if len(client.lastRegister.edKey) == 0 || len(client.lastRegister.authKey) == 0 {
		t.Errorf("keys should not be empty")
	}
}

func TestRegister_ErrorFromClient(t *testing.T) {
	sentinel := errors.New("fail")
	client := &stubClient{registerFn: func(context.Context, string, crypto.KDFParameters, []byte, []byte, crypto.AuthKeyAlgorithm) error {
		return sentinel
	}}
	a := New(client, &stubTokens{}, func([]byte) {})

	if err := a.Register("bob", "pw"); !errors.Is(err, sentinel) {
		t.Fatalf("want %v, got %v", sentinel, err)
	}
}

func TestLogin_SuccessFlow(t *testing.T) {
	password := "pw"
	dataKey, payload, err := makePayload(password)
	if err != nil {
		t.Fatal(err)
	}

	client := &stubClient{
		loginStartFn:  func(context.Context, string, string) (*crypto.LoginPayload, error) { return payload, nil },
		loginFinishFn: func(context.Context, string, string, []byte) (string, error) { return "ACCESS_TOKEN", nil },
	}
	tokens := &stubTokens{}
	var gotDK []byte
	a := New(client, tokens, func(dk []byte) { gotDK = dk })

	if err := a.Login("alice", password); err != nil {
		t.Fatalf("Login unexpected error: %v", err)
	}
	if tokens.lastToken != "ACCESS_TOKEN" {
		t.Errorf("saved token = %q, want ACCESS_TOKEN", tokens.lastToken)
	}
	if !reflect.DeepEqual(gotDK, tokens.lastUnlock) || !reflect.DeepEqual(gotDK, dataKey) {
		t.Errorf("dataKey mismatch between setter/unlock/expected")
	}
}

func TestLogin_ErrorFromStart(t *testing.T) {
	sentinel := errors.New("no user")
	client := &stubClient{loginStartFn: func(context.Context, string, string) (*crypto.LoginPayload, error) { return nil, sentinel }}
	a := New(client, &stubTokens{}, func([]byte) {})

	if err := a.Login("bad", "pw"); !errors.Is(err, sentinel) {
		t.Fatalf("want %v, got %v", sentinel, err)
	}
}

func TestLogin_ErrorFromFinish(t *testing.T) {
	password := "pw"
	_, payload, err := makePayload(password) // valid EncryptedDataKey
	if err != nil {
		t.Fatal(err)
	}

	sentinel := errors.New("bad finish")
	client := &stubClient{
		loginStartFn:  func(context.Context, string, string) (*crypto.LoginPayload, error) { return payload, nil },
		loginFinishFn: func(context.Context, string, string, []byte) (string, error) { return "", sentinel },
	}
	a := New(client, &stubTokens{}, func([]byte) {})

	if err := a.Login("u", password); !errors.Is(err, sentinel) {
		t.Fatalf("want %v, got %v", sentinel, err)
	}
}

func TestLogin_ErrorFromUnlock(t *testing.T) {
	password := "pw"
	_, payload, err := makePayload(password)
	if err != nil {
		t.Fatal(err)
	}

	sentinel := errors.New("locked")
	client := &stubClient{
		loginStartFn:  func(context.Context, string, string) (*crypto.LoginPayload, error) { return payload, nil },
		loginFinishFn: func(context.Context, string, string, []byte) (string, error) { return "TOK", nil },
	}
	tokens := &stubTokens{unlockFn: func([]byte) error { return sentinel }}

	a := New(client, tokens, func([]byte) {})
	if err := a.Login("u", password); !errors.Is(err, sentinel) {
		t.Fatalf("want %v, got %v", sentinel, err)
	}
}

func TestLogin_ErrorFromSaveToken(t *testing.T) {
	password := "pw"
	_, payload, err := makePayload(password)
	if err != nil {
		t.Fatal(err)
	}

	sentinel := errors.New("save failed")
	client := &stubClient{
		loginStartFn:  func(context.Context, string, string) (*crypto.LoginPayload, error) { return payload, nil },
		loginFinishFn: func(context.Context, string, string, []byte) (string, error) { return "TOK", nil },
	}
	tokens := &stubTokens{saveAccessTokenFn: func(context.Context, string) error { return sentinel }}

	a := New(client, tokens, func([]byte) {})
	if err := a.Login("u", password); !errors.Is(err, sentinel) {
		t.Fatalf("want %v, got %v", sentinel, err)
	}
}

func makePayload(password string) ([]byte, *crypto.LoginPayload, error) {
	kdf := crypto.DefaultKDFParameters()
	kdf.Salt = []byte("0123456789abcdef0123456789abcdef") // 32b

	mk, err := crypto.DeriveKey(password, kdf)
	if err != nil {
		return nil, nil, err
	}
	dataKey := []byte("12345678901234567890123456789012") // 32b
	edKey, err := crypto.Encrypt(mk, dataKey)
	if err != nil {
		return nil, nil, err
	}
	return dataKey, &crypto.LoginPayload{
		DeviceId:         "dev",
		KDFParameters:    kdf,
		EncryptedDataKey: edKey,
		AuthKeyAlgorithm: crypto.DefaultAuthKeyAlgorithm(),
		Challenge:        []byte("CHALL"),
	}, nil
}
