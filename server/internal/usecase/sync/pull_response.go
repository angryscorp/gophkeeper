package sync

import "gophkeeper/server/internal/domain"

// PullResponse represents the result of a sync Pull operation.
// It contains a slice of changes, the new cursor, and a flag
// indicating if more data is available.
type PullResponse struct {
	Changes    []domain.Message
	NextCursor int64
	HasMore    bool
}
