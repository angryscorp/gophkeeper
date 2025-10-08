package list

import "context"

type mockRepo struct {
	rows []RawRecord
	err  error
}

func (m mockRepo) GetAll(ctx context.Context) ([]RawRecord, error) {
	return m.rows, m.err
}
