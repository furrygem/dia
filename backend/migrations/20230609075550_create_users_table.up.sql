CREATE TABLE users (
	id UUID PRIMARY KEY,
	username VARCHAR(50) NOT NULL,
	pfp_url VARCHAR(600) NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL,
	updated_at TIMESTAMP WITH TIME ZONE NULL,
	active boolean,
	hashed_password VARCHAR(100)
);


CREATE UNIQUE INDEX IF NOT EXISTS username_index ON users(username)
