-- name: GetChanges :many
SELECT j.server_seq, j.id, j.kind, j.updated_at_unix, j.payload, j.operation_id
FROM journal j
JOIN users u
    ON u.id = j.user_id
WHERE u.username = $1
  AND j.server_seq > $2
ORDER BY j.server_seq
LIMIT $3;
