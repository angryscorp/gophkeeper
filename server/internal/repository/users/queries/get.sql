-- name: Get :one
SELECT
    id, username, kdf_algorithm, kdf_time_cost, kdf_memory_cost, kdf_parallelism, kdf_salt, encrypted_data_key, auth_key, auth_key_algorithm
FROM
    users
WHERE
    username = @username;