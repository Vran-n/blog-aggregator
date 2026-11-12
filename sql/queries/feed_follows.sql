-- name: CreateFeedFollow :one
WITH new_feed_follows AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT 
    new_feed_follows.*,
    users.name AS username,
    feeds.name AS feedname
FROM
    new_feed_follows
    INNER JOIN users ON new_feed_follows.user_id = users.id
    INNER JOIN feeds ON new_feed_follows.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT
    feed_follows.*,
    users.name AS username,
    feeds.name AS feedname,
    feeds.url AS feedurl
FROM
    feed_follows
    INNER JOIN users ON feed_follows.user_id = users.id
    INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE
    feed_follows.user_id = $1;
