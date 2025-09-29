-- name: GetForUpdate :one
SELECT
    c.id AS challenge_id,
    c.challenge,
    c.attempts,
    u.id AS user_id,
    u.auth_key,
    u.auth_key_algorithm
FROM challenges AS c
    JOIN users AS u ON u.id = c.user_id
WHERE u.username = $1
  AND c.device_name = $2
  AND c.used_at IS NULL
  AND c.expires_at > now()
    FOR UPDATE OF c;
