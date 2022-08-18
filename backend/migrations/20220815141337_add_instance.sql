CREATE TABLE instance (
    id INT PRIMARY KEY CHECK (id = 0),
    is_initialized INT NOT NULL CHECK (is_initialized = FALSE OR is_initialized = TRUE)
) STRICT;