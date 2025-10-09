package auth

import (
	"context"
	"crypto/hmac"
	"log"
	"time"

	"gophkeeper/pkg/crypto"
	"gophkeeper/server/internal/domain"
)

const challengeVerificationAttempts = 3

// Auth coordinates registration and login flows.
// It reads users, manages HMAC-based login challenges, and issues access tokens.
type Auth struct {
	users       Users
	challenges  Challenges
	tokenIssuer TokenIssuer
}

// New constructs an Auth service with repositories for users and challenges
// and a token issuer used to mint access tokens after successful login.
func New(
	users Users,
	challenges Challenges,
	tokenIssuer TokenIssuer,
) *Auth {
	return &Auth{
		users:       users,
		challenges:  challenges,
		tokenIssuer: tokenIssuer,
	}
}

// Register creates a new user account in the repository.
// Returns an error if the username is already taken or persistence fails.
func (auth *Auth) Register(ctx context.Context, user domain.User) error {
	log.Printf("Registering user: %s\n", user.Username)
	return auth.users.Add(ctx, user)
}

// LoginStart begins the login flow for the given user/device.
// It loads the user's KDF and encrypted data key, generates a short challenge,
// stores it server-side with expiry, and returns the payload needed by the client
// to derive keys and compute the HMAC response.
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

// LoginFinish completes the login by verifying the client's HMAC over the challenge.
// On success it issues and returns a signed access token; otherwise an error is returned.
func (auth *Auth) LoginFinish(ctx context.Context, username, deviceName string, challenge []byte) (string, error) {
	log.Printf("Finishing login: %s\n", deviceName)

	challengeIsCorrect := false
	err := auth.challenges.GetForUpdate(ctx, username, deviceName, func(info ChallengeInfo) bool {
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
		return "", domain.ErrChallengeFailed
	}

	token, err := auth.tokenIssuer.IssueAccess(username, deviceName)
	if err != nil {
		return "", err
	}

	return token, nil
}
