-- name: UpdateWithSuccess :exec
UPDATE challenges
SET used_at = now()
WHERE id = $1;
