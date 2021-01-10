package db

const postgresSchema = `
	SET TIME ZONE 'UTC';

	CREATE SCHEMA identity;

	CREATE TABLE identity.user (
		user_id SERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);

	CREATE TABLE identity.auth_client (
		auth_client_id SERIAL PRIMARY KEY,
		client_reference TEXT NOT NULL,
		user_id INT REFERENCES identity.user,
		UNIQUE (user_id, client_reference)
	);

	CREATE TYPE token_type AS ENUM ('refresh', 'access');

	CREATE TABLE identity.client_token (
		auth_client_id INT REFERENCES identity.auth_client,
		value TEXT NOT NULL,
		type token_type NOT NULL,
		expired_at TIMESTAMP WITH TIME ZONE NOT NULL,
		PRIMARY KEY (auth_client_id, value, type)
	);
`
