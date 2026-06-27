-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (
        id,
        created_at,
        updated_at,
        user_id,
        feed_id
    ) VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)

SELECT
    i.*,
    f.name AS feed_name,
    u.name AS user_name
FROM inserted_feed_follow AS i
INNER JOIN feeds AS f
    ON i.feed_id = f.id
INNER JOIN users AS u
    ON i.user_id = u.id;

-- name: GetFeedFollowsForUser :many
SELECT
    ff.*,
    f.name AS feed_name,
    u.name AS user_name
FROM feed_follows AS ff
INNER JOIN feeds AS f
    ON ff.feed_id = f.id
INNER JOIN users AS u
    ON ff.user_id = u.id
WHERE u.name = $1;

-- name: DeleteFollow :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;
