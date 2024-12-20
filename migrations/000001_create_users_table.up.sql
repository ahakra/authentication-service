CREATE TABLE IF NOT EXISTS users (
                                     id INTEGER PRIMARY KEY AUTOINCREMENT,
                                     created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     name TEXT NOT NULL,
                                     email TEXT UNIQUE NOT NULL,
                                     password_hash BLOB NOT NULL,
                                     activated BOOLEAN NOT NULL CHECK (activated IN (0, 1)),
    version INTEGER NOT NULL DEFAULT 1
    );
