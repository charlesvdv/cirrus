CREATE TABLE user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role INTEGER NOT NULL,
    FOREIGN KEY(role) REFERENCES role(id)
) STRICT;

CREATE TABLE role (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
) STRICT;

-- add default roles
INSERT INTO role(name) VALUES('User');
INSERT INTO role(name) VALUES('Administrator');