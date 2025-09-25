-- name: GetAccessToken :many
SELECT access_token
FROM tokens
WHERE id = 1;