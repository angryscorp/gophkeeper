-- name: UpdateWithFailure :exec
UPDATE challenges
SET attempts = attempts + 1
WHERE id = $1 AND used_at IS NULL;