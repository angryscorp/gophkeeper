package sync

import "gophkeeper/server/internal/domain"

type PullResponse struct {
	Changes    []domain.Message
	NextCursor int64
	HasMore    bool
}
