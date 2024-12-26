CREATE TABLE IF NOT EXISTS permissions (
                                           id integer PRIMARY KEY AUTOINCREMENT,
                                           permission text unique NOT NULL
);

CREATE TABLE IF NOT EXISTS users_permissions (
                                                 user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
                                                 permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
                                                 PRIMARY KEY (user_id, permission_id)
    );
INSERT INTO permissions (permission)
values('permissions:read'),
      ('permissions:write')
