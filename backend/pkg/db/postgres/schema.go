package postgres

import (
	"database/sql"
)

const sqlSchema = `
	SET TIMEZONE = 'UTC';
	CREATE SCHEMA cirrus;

	CREATE TABLE cirrus.inode (
		id UUID NOT NULL,
		parent_id UUID,
		name TEXT NOT NULL,
		created_time TIMESTAMP,
		PRIMARY KEY (id)
	)
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
