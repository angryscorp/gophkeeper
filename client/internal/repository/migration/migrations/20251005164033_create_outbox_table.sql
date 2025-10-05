-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS outbox (
    operation_id     TEXT PRIMARY KEY,
    record_id        TEXT NOT NULL,
    kind             INTEGER NOT NULL,
    updated_at_unix  INTEGER NOT NULL,  -- ms, UTC
    payload          BLOB NOT NULL,     -- = nonce||ciphertext||tag
    created_at_unix  INTEGER NOT NULL   -- sorting
);

CREATE INDEX IF NOT EXISTS idx_outbox_created ON outbox(created_at_unix);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS outbox;
-- +goose StatementEnd
