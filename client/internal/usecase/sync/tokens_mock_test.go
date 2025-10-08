package sync

import "context"

type mockTokens struct {
	ready bool
	token string
	err   error
}

func (m *mockTokens) GetAccessToken(ctx context.Context) (string, error) {
	return m.token, m.err
}
func (m *mockTokens) Ready() bool { return m.ready }
