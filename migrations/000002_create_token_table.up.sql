CREATE TABLE IF NOT EXISTS tokens (
        hash BLOB NOT NULL PRIMARY KEY,
        user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
        expiry timestamp(0)  NOT NULL,
        scope text NOT NULL
                                                              );