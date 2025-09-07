package repository

// Database: SQLite (mattn/go-sqlite3) + SQLCipher

// Unlock
// - User enters master password
// - MasterKey = Argon2id(password, salt, t/m/p)
// - Set key for SQLCipher as raw key: PRAGMA key = "x'<HEX_32_BYTES>'"
// - After PRAGMA key regular SQL queries, database is "transparently" decrypted

// Change password
// - New MasterKey (Argon2id with new salt/parameters)
// - PRAGMA rekey = "x'<NEW_HEX_32_BYTES>'";
