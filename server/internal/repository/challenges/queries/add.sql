-- name: Add :exec
INSERT INTO
    challenges (user_id, id, device_name, challenge, expires_at)
VALUES ($1, $2, $3, $4 ,$5);
