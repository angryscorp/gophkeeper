-- name: Add :exec
INSERT INTO users (
    id, username,
    kdf_algorithm, kdf_time_cost, kdf_memory_cost, kdf_parallelism, kdf_salt,
    encrypted_data_key,
    auth_key, auth_key_algorithm
) VALUES (
    $1, $2,
    $3, $4, $5, $6, $7,
    $8,
    $9, $10
);