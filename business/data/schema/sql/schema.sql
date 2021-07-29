-- Version: 1.0
-- Description: Create table users
CREATE TABLE auth_user (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    passw TEXT NOT NULL,
    created DATETIME NOT NULL
);

CREATE INDEX auth_user1 ON auth_user(user_id);
CREATE UNIQUE INDEX auth_user_id_UNIQUE ON auth_user(email ASC);

