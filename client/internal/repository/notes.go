package repository

// Unlock
// - User enters master password
// - MasterKey = Argon2id(password, salt, t/m/p)
// - DBKey := HKDF(MK, "sqlcipher", 32)
// - Set key for SQLCipher as raw key: PRAGMA key = "x'<HEX_32_BYTES>'"
// - After PRAGMA key regular SQL queries, database is "transparently" decrypted
