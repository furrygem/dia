CREATE TABLE IF NOT EXISTS publickeys (
	fingerprint bytea PRIMARY KEY,
	publickey VARCHAR
);

CREATE INDEX iF NOT EXISTS publickey_index ON publickeys(publickey);
