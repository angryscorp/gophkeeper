-- +goose Up
-- +goose StatementBegin
CREATE TABLE journal (
    user_id         UUID    NOT NULL,
    server_seq      BIGSERIAL PRIMARY KEY,  -- порядок Pull
    id              UUID    NOT NULL,       -- record_id
    kind            INT     NOT NULL,
    updated_at_unix BIGINT  NOT NULL,       -- LWW (ms, UTC)
    payload         BYTEA   NOT NULL,       -- nonce||cipher||tag
    operation_id    UUID    NOT NULL        -- idempotency
);

CREATE UNIQUE INDEX ON journal (user_id, operation_id);
CREATE INDEX        ON journal (user_id, id, updated_at_unix DESC, server_seq DESC);
CREATE INDEX        ON journal (user_id, server_seq);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS journal;
-- +goose StatementEnd
