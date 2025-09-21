package usecase

// Sign-Up
// - The user enters a username and master password
// - The client calculates Master Key = Argon2id(password, salt, params)
// - A random Data Key is generated → it encrypts all records
// - The Data Key is encrypted by the Master Key → Encrypted Data Key
// - AuthKey = HKDF(DataKey, "auth") is derived from the Data Key
// - The following is sent to the server: login, salt+KDF parameters, Encrypted Data Key, AuthKey
// - The server saves them, but does not know the password and does not store the keys in plain text

// Sign-In
// - The client receives from the server by username: salt + KDF parameters, Encrypted Data Key, and challenge
// - The client calculates the Master Key by password → decrypts the Encrypted Data Key → receives the Data Key
// - The client derives AuthKey = HKDF(DataKey, "auth") and signs the challenge: HMAC(AuthKey, challenge)
// - The server checks via its AuthKey
// - If everything is ok → the server issues an access-token + refresh-token

// Password change
// - The client decrypts the Data Key with the old Master Key
// - Calculates the new Master Key from the new password
// - Re-encrypts the Data Key → new Encrypted Data Key
// - Sends it and the new KDF parameters to the server
// - The AuthKey does not change because the Data Key remains the same → synchronization is not broken

// Sync:
// - Push → Pull
// - Conflicts: Last-Write-Wins
// - All sensitive data (logins/passwords/notes/etc) are encrypted on the client
