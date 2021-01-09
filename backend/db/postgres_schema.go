package db

const postgresSchema = `
	SET TIME ZONE 'UTC';

	CREATE SCHEMA identity;

	CREATE TABLE identity.user (
		user_id SERIAL PRIMARY KEY,
		email TEXT UNIQUE,
		password TEXT
	);
`
