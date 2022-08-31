CREATE TABLE session (
    token TEXT NOT NULL UNIQUE,
    user INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
    expired_at INTEGER NOT NULL -- stored as a UNIX epoch
) STRICT;