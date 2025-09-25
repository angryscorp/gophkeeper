package auth

import (
	"context"
	"crypto/hmac"
	"errors"
	"gophkeeper/pkg/crypto"
	"gophkeeper/server/internal/domain"
	"gophkeeper/server/internal/repository/challenges"
	"gophkeeper/server/internal/repository/users"
	"log"
	"time"
)

const challengeVerificationAttempts = 3

type TokenIssuer interface {
	IssueAccess(userID, deviceID string) (string, error)
}
type Auth struct {
	users       users.Users
	challenges  challenges.Challenges
	tokenIssuer TokenIssuer
}

func New(
	users users.Users,
	challenges challenges.Challenges,
	tokenIssuer TokenIssuer,
) *Auth {
	return &Auth{
		users:       users,
		challenges:  challenges,
		tokenIssuer: tokenIssuer,
	}
}

func (auth *Auth) Register(ctx context.Context, user domain.User) error {
	log.Printf("Registering user: %s\n", user.Username)
	return auth.users.Add(ctx, user)
}

func (auth *Auth) LoginStart(ctx context.Context, username, deviceId string) (crypto.LoginPayload, error) {
	log.Printf("Starting login for user: %s\n", username)
	resp, err := auth.users.Get(ctx, username)
	if err != nil {
		return crypto.LoginPayload{}, err
	}

	ch := crypto.RandBytes(8)
	err = auth.challenges.Add(ctx, resp.ID, deviceId, ch, time.Now().Add(time.Minute))
	if err != nil {
		return crypto.LoginPayload{}, err
	}

	return crypto.LoginPayload{
		DeviceId:         deviceId,
		KDFParameters:    resp.KDFParameters,
		EncryptedDataKey: resp.EncryptedDataKey,
		AuthKeyAlgorithm: resp.AuthKeyAlgorithm,
		Challenge:        ch,
	}, nil
}

func (auth *Auth) LoginFinish(ctx context.Context, username, deviceName string, challenge []byte) (string, error) {
	log.Printf("Finishing login: %s\n", deviceName)

	challengeIsCorrect := false
	err := auth.challenges.GetForUpdate(ctx, username, deviceName, func(info challenges.ChallengeInfo) bool {
		if info.Attempts >= challengeVerificationAttempts {
			return false
		}
		expected := crypto.SignChallenge(info.AuthKey, info.Challenge, info.AuthKeyAlgorithm)
		challengeIsCorrect = hmac.Equal(expected, challenge)
		return challengeIsCorrect
	})
	if err != nil {
		return "", err
	}

	if !challengeIsCorrect {
		return "", errors.New("invalid challenge")
	}

	token, err := auth.tokenIssuer.IssueAccess(username, deviceName)
	if err != nil {
		return "", err
	}

	return token, nil
}
