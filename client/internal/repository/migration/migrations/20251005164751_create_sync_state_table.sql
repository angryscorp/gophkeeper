-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sync_state (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    last_server_seq INTEGER NOT NULL DEFAULT 0
);
INSERT OR IGNORE INTO sync_state(id, last_server_seq) VALUES (1, 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sync_state
-- +goose StatementEnd
