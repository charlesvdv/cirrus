CREATE TABLE instance (
    id INTEGER PRIMARY KEY CHECK (id = 0),
    is_initialized INTEGER NOT NULL CHECK (is_initialized = FALSE OR is_initialized = TRUE)
) STRICT;