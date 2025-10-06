-- name: UpsertRecord :exec
INSERT INTO records (id, kind, updated_at_unix, payload)
VALUES (?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    kind            = excluded.kind,
    updated_at_unix = excluded.updated_at_unix,
    payload         = excluded.payload;

-- name: GetRecordMeta :one
SELECT updated_at_unix FROM records WHERE id = ?;
