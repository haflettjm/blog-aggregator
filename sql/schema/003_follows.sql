-- +goose up
CREATE TABLE feed_follows (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	feed_id INT NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
	UNIQUE(user_id, feed_id)
);


-- +goose down
DROP TABLE feeds;
