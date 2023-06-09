CREATE TABLE users (
	id UUID PRIMARY KEY,
	username VARCHAR(50) NOT NULL,
	pfpurl VARCHAR(600) NULL,
	createdat TIMESTAMP NOT NULL,
	updatedat TIMESTAMP NULL,
	active boolean
);


CREATE UNIQUE INDEX IF NOT EXISTS username_index ON users(username)
