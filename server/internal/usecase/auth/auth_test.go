package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"gophkeeper/pkg/crypto"
	"gophkeeper/server/internal/domain"

	"github.com/google/uuid"
)

func TestRegister_SuccessAndError(t *testing.T) {
	var added domain.User
	usersOK := stubUsers{
		addFn: func(ctx context.Context, u domain.User) error {
			added = u
			return nil
		},
	}
	a := New(usersOK, stubChallenges{}, stubIssuer{})
	u := domain.User{Username: "alice"}
	if err := a.Register(context.Background(), u); err != nil {
		t.Fatalf("Register() unexpected error: %v", err)
	}
	if added.Username != "alice" {
		t.Fatalf("Register() did not pass user to repo")
	}

	sentinel := errors.New("duplicate")
	usersErr := stubUsers{
		addFn: func(ctx context.Context, u domain.User) error { return sentinel },
	}
	a = New(usersErr, stubChallenges{}, stubIssuer{})
	if err := a.Register(context.Background(), u); !errors.Is(err, sentinel) {
		t.Fatalf("Register() want %v, got %v", sentinel, err)
	}
}

func TestLoginStart_Success(t *testing.T) {
	user := &domain.User{
		ID:               uuid.New(),
		Username:         "bob",
		KDFParameters:    crypto.DefaultKDFParameters(),
		EncryptedDataKey: []byte{1, 2, 3},
		AuthKeyAlgorithm: crypto.DefaultAuthKeyAlgorithm(),
	}

	users := stubUsers{
		getFn: func(ctx context.Context, username string) (*domain.User, error) {
			if username != "bob" {
				t.Fatalf("unexpected username: %s", username)
			}
			return user, nil
		},
	}

	var captured struct {
		userID     uuid.UUID
		deviceName string
		challenge  []byte
		expiresAt  time.Time
	}

	chal := stubChallenges{
		addFn: func(ctx context.Context, userId uuid.UUID, deviceName string, challenge []byte, expiresAt time.Time) error {
			captured.userID = userId
			captured.deviceName = deviceName
			captured.challenge = challenge
			captured.expiresAt = expiresAt
			return nil
		},
	}

	a := New(users, chal, stubIssuer{})
	resp, err := a.LoginStart(context.Background(), "bob", "dev-1")
	if err != nil {
		t.Fatalf("LoginStart() unexpected error: %v", err)
	}

	if resp.DeviceId != "dev-1" {
		t.Fatalf("DeviceId mismatch: %q", resp.DeviceId)
	}
	if len(resp.Challenge) != 8 {
		t.Fatalf("challenge length: want 8, got %d", len(resp.Challenge))
	}
	if string(resp.EncryptedDataKey) != string(user.EncryptedDataKey) {
		t.Fatalf("EncryptedDataKey mismatch")
	}
	if captured.userID != user.ID || captured.deviceName != "dev-1" {
		t.Fatalf("challenge Add() not called with expected args")
	}
	if time.Until(captured.expiresAt) <= 0 {
		t.Fatalf("expiresAt should be in the future")
	}
}

func TestLoginFinish_Success(t *testing.T) {
	authKey := []byte("0123456789abcdef0123456789abcdef")
	ch := []byte("abcdefgh")
	algo := crypto.AuthKeyAlgorithmHMACSHA256

	challenges := stubChallenges{
		getForUpdateFn: func(ctx context.Context, username, device string, validator func(ChallengeInfo) bool) error {
			info := ChallengeInfo{
				Challenge:        ch,
				Attempts:         0,
				AuthKey:          authKey,
				AuthKeyAlgorithm: algo,
			}
			_ = validator(info) // should be true for correct HMAC
			return nil
		},
	}

	issuer := stubIssuer{
		issueFn: func(userID, deviceID string) (string, error) {
			if userID != "bob" || deviceID != "dev-1" {
				t.Fatalf("IssueAccess args mismatch")
			}
			return "ACCESS_TOKEN", nil
		},
	}

	a := New(stubUsers{}, challenges, issuer)
	resp := crypto.SignChallenge(authKey, ch, algo)
	token, err := a.LoginFinish(context.Background(), "bob", "dev-1", resp)
	if err != nil {
		t.Fatalf("LoginFinish() unexpected error: %v", err)
	}
	if token != "ACCESS_TOKEN" {
		t.Fatalf("token mismatch: %q", token)
	}
}

func TestLoginFinish_WrongHMAC(t *testing.T) {
	authKey := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	ch := []byte("abcdefgh")
	algo := crypto.AuthKeyAlgorithmHMACSHA256

	challenges := stubChallenges{
		getForUpdateFn: func(ctx context.Context, username, device string, validator func(ChallengeInfo) bool) error {
			info := ChallengeInfo{
				Challenge:        ch,
				Attempts:         0,
				AuthKey:          authKey,
				AuthKeyAlgorithm: algo,
			}
			// client will pass WRONG response; validator will return false
			_ = validator(info)
			return nil
		},
	}

	a := New(stubUsers{}, challenges, stubIssuer{})
	// wrong response (different key)
	wrong := crypto.SignChallenge([]byte("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"), ch, algo)

	_, err := a.LoginFinish(context.Background(), "u", "d", wrong)
	if !errors.Is(err, domain.ErrChallengeFailed) {
		t.Fatalf("want ErrChallengeFailed, got %v", err)
	}
}

func TestLoginFinish_GetForUpdateError(t *testing.T) {
	sentinel := errors.New("storage down")
	challenges := stubChallenges{
		getForUpdateFn: func(ctx context.Context, username, device string, validator func(ChallengeInfo) bool) error {
			return sentinel
		},
	}
	a := New(stubUsers{}, challenges, stubIssuer{})
	_, err := a.LoginFinish(context.Background(), "u", "d", []byte("x"))
	if !errors.Is(err, sentinel) {
		t.Fatalf("want %v, got %v", sentinel, err)
	}
}

func TestLoginFinish_IssuerError(t *testing.T) {
	authKey := []byte("0123456789abcdef0123456789abcdef")
	ch := []byte("abcdefgh")
	algo := crypto.AuthKeyAlgorithmHMACSHA256

	challenges := stubChallenges{
		getForUpdateFn: func(ctx context.Context, username, device string, validator func(ChallengeInfo) bool) error {
			_ = validator(ChallengeInfo{
				Challenge:        ch,
				Attempts:         0,
				AuthKey:          authKey,
				AuthKeyAlgorithm: algo,
			})
			return nil
		},
	}
	issuer := stubIssuer{
		issueFn: func(userID, deviceID string) (string, error) { return "", errors.New("sign failed") },
	}

	a := New(stubUsers{}, challenges, issuer)
	resp := crypto.SignChallenge(authKey, ch, algo)
	_, err := a.LoginFinish(context.Background(), "u", "d", resp)
	if err == nil || err.Error() != "sign failed" {
		t.Fatalf("expected issuer error, got %v", err)
	}
}
