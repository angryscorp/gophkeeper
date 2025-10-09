-- name: ListOutboxBatch :many
SELECT operation_id, record_id, kind, updated_at_unix, payload
FROM outbox
ORDER BY created_at_unix
LIMIT ?;

-- name: DeleteOutbox :exec
DELETE FROM outbox WHERE operation_id = ?;