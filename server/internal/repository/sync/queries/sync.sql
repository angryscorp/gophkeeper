-- name: GetChanges :many
SELECT j.server_seq, j.id, j.kind, j.updated_at_unix, j.payload, j.operation_id
FROM journal j
JOIN users u
    ON u.id = j.user_id
WHERE u.username = $1
  AND j.server_seq > $2
ORDER BY j.server_seq
LIMIT $3;

-- name: InsertChange :exec
WITH u AS (
    SELECT id FROM users WHERE username = $1
)
INSERT INTO journal(user_id, id, kind, updated_at_unix, payload, operation_id)
SELECT u.id, $2, $3, $4, $5, $6
FROM u
ON CONFLICT (user_id, operation_id) DO NOTHING;