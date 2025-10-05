-- name: AddRecord :exec
INSERT INTO records (id, kind, updated_at_unix, payload)
VALUES (?, ?, ?, ?);

-- name: GetRecords :many
SELECT id, kind, payload
FROM records;

