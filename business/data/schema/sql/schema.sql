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

-- Version: 1.1
-- Description: Create table contests
CREATE TABLE contests (
    contest_id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT NOT NULL,
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    created_on DATETIME NOT NULL,
    updated_on DATETIME NOT NULL,
    updated_by TEXT NOT NULL
);

CREATE INDEX contests1 ON contests(contest_id);
CREATE UNIQUE INDEX slug_UNIQUE ON contests(slug ASC);

-- Version: 1.2
-- Description: Create table photos
CREATE TABLE photos (
    photo_id TEXT NOT NULL, --uuid
    owner_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    created_on DATETIME NOT NULL,
    updated_on DATETIME NOT NULL,
    updated_by TEXT NOT NULL,
    FOREIGN KEY(owner_id) REFERENCES auth_user(user_id)
);

CREATE INDEX photosp_hoto_id_ndx ON photos(photo_id);
CREATE INDEX photos_user_id_ndx ON photos(owner_id);


