-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS records (
    id              TEXT PRIMARY KEY,  -- UUID as String
    kind            INTEGER NOT NULL,  -- 1=BankCard, 2=Credentials, 3=Text, 4=Binary
    updated_at_unix INTEGER NOT NULL,  -- ms (UTC)
    payload         BLOB    NOT NULL   -- = nonce(12) || ciphertext || tag(16)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS records;
-- +goose StatementEnd