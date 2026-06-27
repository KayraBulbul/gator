-- name: CreateFeed :one
INSERT INTO feeds (
    id,
    created_at,
    updated_at,
    name,
    url,
    user_id,
    last_fetched_at
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetFeeds :many
SELECT
    f.name AS feedname,
    f.url,
    u.name AS username
FROM feeds AS f
INNER JOIN users AS u
    ON f.user_id = u.id;

-- name: GetFeedByURL :one
SELECT
    f.id,
    f.name AS feed_name,
    u.name AS user_name
FROM feeds AS f
INNER JOIN users AS u
    ON f.user_id = u.id
WHERE f.url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
