package repository

// Database: SQLite (mattn/go-sqlite3) + SQLCipher

// Unlock
// - User enters master password
// - MasterKey = Argon2id(password, salt, t/m/p)
// - DBKey := HKDF(MK, "sqlcipher", 32)
// - Set key for SQLCipher as raw key: PRAGMA key = "x'<HEX_32_BYTES>'"
// - After PRAGMA key regular SQL queries, database is "transparently" decrypted

// Change password
// - MasterKey := Argon2id(masterPassword, salt, params)
// - DBKey := HKDF(MK, "sqlcipher", 32)
// - PRAGMA rekey = "x'<NEW_HEX_32_BYTES>'";
