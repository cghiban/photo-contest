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
-- Description: Create table contest_entries
CREATE TABLE contest_entries (
    entry_id INTEGER PRIMARY KEY AUTOINCREMENT,
    contest_id INTEGER NOT NULL DEFAULT 0,
    photo_id TEXT NOT NULL, --uuid
    status TEXT NOT NULL, -- regex='^(active|eliminated|withdrawn|flagged)$')
    created_on DATETIME NOT NULL,
    updated_on DATETIME NOT NULL,
    updated_by TEXT NOT NULL,
    FOREIGN KEY(photo_id) REFERENCES photos(photo_id),
    FOREIGN KEY(contest_id) REFERENCES contests(contest_id)
);

CREATE INDEX cp_id_ndx ON contest_entries(entry_id);
CREATE INDEX cp_contest_id_ndx ON contest_entries(contest_id);
CREATE INDEX cp_photo_id_ndx ON contest_entries(photo_id);
CREATE INDEX cp_status_ndx ON contest_entries(status);

-- Version: 1.5
-- Description: Create table contest_entry_votes
CREATE TABLE contest_entry_votes (
    v_id INTEGER PRIMARY KEY AUTOINCREMENT,
    v_entry_id INTEGER NOT NULL DEFAULT 0,
    v_contest_id INTEGER NOT NULL DEFAULT 0,
    v_photo_id TEXT NOT NULL, --uuid
    v_user_id INTEGER NOT NULL,
    v_score INTEGER NOT NULL DEFAULT 1,
    v_created_on DATETIME NOT NULL,
    FOREIGN KEY(v_entry_id) REFERENCES contest_entries(entry_id),
    FOREIGN KEY(v_photo_id) REFERENCES photos(photo_id),
    FOREIGN KEY(v_user_id) REFERENCES auth_user(user_id),
    FOREIGN KEY(v_contest_id) REFERENCES contests(contest_id)
);

CREATE INDEX cev_id_ndx ON contest_entry_votes(v_id);
CREATE INDEX cev_entry_id_ndx ON contest_entry_votes(v_entry_id);
CREATE INDEX cev_contest_id_ndx ON contest_entry_votes(v_contest_id);
CREATE INDEX cev_photo_id_ndx ON contest_entry_votes(v_photo_id);
CREATE UNIQUE INDEX cev_unique_ndx ON contest_entry_votes(v_contest_id, v_photo_id, v_user_id);
CREATE UNIQUE INDEX cev_unique_entry_ndx ON contest_entry_votes(v_entry_id, v_user_id);

-- Version: 1.6
-- Description: Add more information for authorized users

ALTER TABLE auth_user ADD COLUMN street TEXT NOT NULL DEFAULT "";
ALTER TABLE auth_user ADD COLUMN city TEXT NOT NULL DEFAULT "";
ALTER TABLE auth_user ADD COLUMN state TEXT NOT NULL DEFAULT "OO";
ALTER TABLE auth_user ADD COLUMN zip TEXT NOT NULL DEFAULT "";
ALTER TABLE auth_user ADD COLUMN phone TEXT NOT NULL DEFAULT "";
ALTER TABLE auth_user ADD COLUMN age INTEGER NOT NULL DEFAULT 0;
ALTER TABLE auth_user ADD COLUMN gender TEXT NOT NULL DEFAULT "-";
ALTER TABLE auth_user ADD COLUMN ethnicity TEXT NOT NULL DEFAULT "pn";
ALTER TABLE auth_user ADD COLUMN other_ethnicity TEXT NOT NULL DEFAULT "";

-- Version: 1.7
-- Description: Add information about subject to contest_entries

ALTER TABLE contest_entries ADD COLUMN sname TEXT NOT NULL DEFAULT "";
ALTER TABLE contest_entries ADD COLUMN sage TEXT NOT NULL DEFAULT "0";
ALTER TABLE contest_entries ADD COLUMN scountry TEXT NOT NULL DEFAULT "US";
ALTER TABLE contest_entries ADD COLUMN sorigin TEXT NOT NULL DEFAULT "";
ALTER TABLE contest_entries ADD COLUMN location TEXT NOT NULL DEFAULT "";
ALTER TABLE contest_entries ADD COLUMN sbiography TEXT NOT NULL DEFAULT "";
ALTER TABLE contest_entries ADD COLUMN release_mime_type TEXT NOT NULL DEFAULT "application/pdf";

-- Version: 1.8
-- Description: Create table reset_password_email
CREATE TABLE reset_password_email (
    reset_id TEXT NOT NULL PRIMARY KEY, --uuid
    user_id INTEGER NOT NULL,
    active BOOLEAN NOT NULL DEFAULT 1,
    created_on DATETIME NOT NULL,
    updated_on DATETIME NOT NULL,
    updated_by TEXT NOT NULL,
    FOREIGN KEY(user_id) REFERENCES auth_user(user_id)
);