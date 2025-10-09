-- name: EnqueueOutbox :exec
INSERT INTO outbox (operation_id, record_id, kind, updated_at_unix, payload, created_at_unix)
VALUES (?, ?, ?, ?, ?, ?);
