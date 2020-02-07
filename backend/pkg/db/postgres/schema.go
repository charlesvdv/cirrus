package postgres

import (
	"database/sql"
)

const sqlSchema = `
	SET TIMEZONE = 'UTC';
	CREATE SCHEMA cirrus;

	CREATE TYPE metadata_type AS ENUM ('file', 'directory');

	CREATE TABLE cirrus.metadata (
		id UUID NOT NULL,
		parent_id UUID NOT NULL,
		name TEXT NOT NULL,
		created_time TIMESTAMP NOT NULL,
		type metadata_type,
		PRIMARY KEY (id),
		UNIQUE (parent_id, name)
	);
`

func setupSchema(db *sql.DB) error {
	schemaInstalled, err := isSchemaInstalled(db)
	if err != nil {
		return err
	}
	if !schemaInstalled {
		if err := installSchema(db); err != nil {
			return err
		}
	}
	return nil
}

func isSchemaInstalled(db *sql.DB) (bool, error) {
	rows, err := db.Query("SELECT 1 FROM information_schema.schemata WHERE schema_name = $1", "cirrus")
	if err != nil {
		return false, err
	}
	return rows.Next(), nil
}

func installSchema(db *sql.DB) error {
	_, err := db.Exec(sqlSchema)
	if err != nil {
		return err
	}
	return nil
}
