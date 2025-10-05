-- name: GetCursor :one
SELECT last_server_seq FROM sync_state WHERE id = 1;

-- name: SetCursor :exec
UPDATE sync_state SET last_server_seq = ? WHERE id = 1;