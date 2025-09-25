-- name: SaveAccessToken :exec
INSERT INTO tokens (id, access_token)
VALUES (1, @token)
ON CONFLICT(id) DO UPDATE
    SET access_token = excluded.access_token,
        created_at   = excluded.created_at;