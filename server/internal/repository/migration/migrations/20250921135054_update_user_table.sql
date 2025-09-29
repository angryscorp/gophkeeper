-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN kdf_algorithm VARCHAR(64) NOT NULL;
ALTER TABLE users ADD COLUMN kdf_time_cost INTEGER NOT NULL;
ALTER TABLE users ADD COLUMN kdf_memory_cost INTEGER NOT NULL;
ALTER TABLE users ADD COLUMN kdf_parallelism INTEGER NOT NULL;
ALTER TABLE users ADD COLUMN kdf_salt BYTEA NOT NULL;
ALTER TABLE users ADD COLUMN encrypted_data_key BYTEA NOT NULL;
ALTER TABLE users ADD COLUMN auth_key BYTEA NOT NULL;
ALTER TABLE users ADD COLUMN auth_key_algorithm VARCHAR(64) NOT NULL;
ALTER TABLE users ADD COLUMN created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();
ALTER TABLE users ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN kdf_algorithm;
ALTER TABLE users DROP COLUMN kdf_time_cost;
ALTER TABLE users DROP COLUMN kdf_memory_cost;
ALTER TABLE users DROP COLUMN kdf_parallelism;
ALTER TABLE users DROP COLUMN kdf_salt;
ALTER TABLE users DROP COLUMN encrypted_data_key;
ALTER TABLE users DROP COLUMN auth_key;
ALTER TABLE users DROP COLUMN auth_key_algorithm;
ALTER TABLE users DROP COLUMN created_at;
ALTER TABLE users DROP COLUMN updated_at;
-- +goose StatementEnd
