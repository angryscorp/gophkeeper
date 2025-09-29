-- name: GetAccessToken :one
SELECT access_token
FROM tokens
WHERE id = 1;