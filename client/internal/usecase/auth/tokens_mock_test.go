package auth

import "context"

type stubTokens struct {
	unlockFn          func(dataKey []byte) error
	saveAccessTokenFn func(ctx context.Context, token string) error

	lastUnlock []byte
	lastToken  string
}

func (s *stubTokens) Unlock(dataKey []byte) error {
	s.lastUnlock = dataKey
	if s.unlockFn != nil {
		return s.unlockFn(dataKey)
	}
	return nil
}

func (s *stubTokens) SaveAccessToken(ctx context.Context, token string) error {
	s.lastToken = token
	if s.saveAccessTokenFn != nil {
		return s.saveAccessTokenFn(ctx, token)
	}
	return nil
}
