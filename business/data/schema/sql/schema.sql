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
    photo_id TEXT NOT NULL PRIMARY KEY, --uuid
    owner_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    deleted BOOLEAN NOT NULL DEFAULT 0,
    created_on DATETIME NOT NULL,
    updated_on DATETIME NOT NULL,
    updated_by TEXT NOT NULL,
    FOREIGN KEY(owner_id) REFERENCES auth_user(user_id)
);

CREATE INDEX photos_photo_id_ndx ON photos(photo_id);
CREATE INDEX photos_user_id_ndx ON photos(owner_id);

-- Version: 1.3
-- Description: Create table photo_files
CREATE TABLE photo_files (
    file_id TEXT NOT NULL PRIMARY KEY, --uuid
    photo_id TEXT NOT NULL, --uuid
    filepath TEXT NOT NULL,
    size TEXT NOT NULL, -- regex='^(thumb|small|medium|large|original|custom)$')
    w INTEGER NOT NULL DEFAULT 0,
    h INTEGER NOT NULL DEFAULT 0,
    created_on DATETIME NOT NULL,
    updated_on DATETIME NOT NULL,
    updated_by TEXT NOT NULL,
    FOREIGN KEY(photo_id) REFERENCES photos(photo_id)
);

CREATE INDEX photo_files_file_id_ndx ON photo_files(file_id);
CREATE INDEX photo_files_photo_id_ndx ON photo_files(photo_id);

-- Version: 1.4
-- Description: Create table contest_photos
CREATE TABLE contest_photos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    contest_id INTEGER NOT NULL DEFAULT 0,
    photo_id TEXT NOT NULL, --uuid
    filepath TEXT NOT NULL,
    status TEXT NOT NULL, -- regex='^(active|eliminated|withdrawn|flagged)$')
    created_on DATETIME NOT NULL,
    updated_on DATETIME NOT NULL,
    updated_by TEXT NOT NULL,
    FOREIGN KEY(photo_id) REFERENCES photos(photo_id),
    FOREIGN KEY(contest_id) REFERENCES contests(contest_id)
);

CREATE INDEX cp_id_ndx ON contest_photos(id);
CREATE INDEX cp_contest_id_ndx ON contest_photos(contest_id);
CREATE INDEX cp_photo_id_ndx ON contest_photos(photo_id);
CREATE INDEX cp_status_ndx ON contest_photos(status);


