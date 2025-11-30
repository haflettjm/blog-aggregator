-- +goose up
CREATE TABLE feeds (
	id SERIAL PRIMARY KEY,
	url TEXT NOT NULL UNIQUE,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	last_updated TIMESTAMP NOT NULL DEFAULT NOW(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);


-- +goose down
DROP TABLE feeds;
