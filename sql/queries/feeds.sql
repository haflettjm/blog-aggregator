-- name: CreateFeed :one
INSERT INTO feeds (url, name, description, created_at, last_updated, user_id)
VALUES ($1, $2, $3, NOW(), NOW(), $4)
RETURNING id;

-- name: UpdateFeed :one
UPDATE feeds
SET url = $2, name = $3, description = $4, last_updated = NOW()
WHERE id = $1
RETURNING id;

-- name: DeleteFeed :one
DELETE FROM feeds
WHERE id = $1
RETURNING id;

-- name: ListFeeds :many
SELECT * FROM feeds WHERE id = $1;

-- name: GetFeedsByUserId :many
SELECT * FROM feeds WHERE user_id = $1;

-- name: GetFeedsByName :many
SELECT * FROM feeds WHERE name = $1;

-- name: GetFeedsByURL :many
SELECT * FROM feeds WHERE url = $1;

-- name: DeleteFeeds :exec
DELETE FROM feeds;
