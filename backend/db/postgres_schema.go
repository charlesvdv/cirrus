package db

const postgresSchema = `
	CREATE SCHEMA usermgt;

	CREATE TABLE usermgt.user (
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE,
		password TEXT
	);
`
