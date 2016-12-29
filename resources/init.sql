--
-- Coalesce Database Initialization Query
--
-- Text encoding used: UTF-8
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: categories
CREATE TABLE IF NOT EXISTS categories (
    categoryid INTEGER PRIMARY KEY ASC AUTOINCREMENT
                       UNIQUE
                       NOT NULL,
    name       TEXT    UNIQUE
                       NOT NULL
);


-- Table: comments
CREATE TABLE IF NOT EXISTS comments (
    commentid INTEGER PRIMARY KEY ASC AUTOINCREMENT
                      UNIQUE
                      NOT NULL,
    postid    INTEGER REFERENCES posts (postid) 
                      NOT NULL,
    parentid  INTEGER REFERENCES comments (commentid),
    userid            REFERENCES users (userid) 
                      NOT NULL,
    body      TEXT    NOT NULL,
    bodyHTML  TEXT    NOT NULL,
    posted    INTEGER NOT NULL,
    updated   INTEGER NOT NULL
);


-- Table: images
CREATE TABLE IF NOT EXISTS images (
    imageid   INTEGER PRIMARY KEY ASC AUTOINCREMENT
                      NOT NULL
                      UNIQUE,
    userid    INTEGER REFERENCES users (userid) 
                      NOT NULL,
    md5       TEXT    NOT NULL,
    thumb_md5 TEXT    NOT NULL,
    filename  TEXT    NOT NULL
);


-- Table: posts
CREATE TABLE IF NOT EXISTS posts (
    postid     INTEGER PRIMARY KEY ASC AUTOINCREMENT
                       UNIQUE
                       NOT NULL,
    userid     INTEGER REFERENCES users (userid) 
                       NOT NULL,
    title      TEXT    NOT NULL,
    body       TEXT,
    bodyHTML   TEXT,
--    categoryid INTEGER REFERENCES categories (categoryid),
    posted     INTEGER NOT NULL,
    updated    INTEGER NOT NULL
);


-- Table: users
CREATE TABLE IF NOT EXISTS users (
    userid    INTEGER PRIMARY KEY ASC AUTOINCREMENT
                      NOT NULL
                      UNIQUE,
    username  TEXT    UNIQUE
                      NOT NULL,
    token     TEXT    NOT NULL,
    privlevel INTEGER NOT NULL
);

-- Table: errors
CREATE TABLE IF NOT EXISTS errors (
    errorid   INTEGER PRIMARY KEY ASC AUTOINCREMENT
                      NOT NULL
                      UNIQUE,
    errortext TEXT    NOT NULL
);

COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
