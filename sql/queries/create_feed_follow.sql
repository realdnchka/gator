-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *, 
(SELECT users.name FROM users WHERE users.id = $4) as user_name, 
(SELECT feeds.name FROM feeds WHERE feeds.id = $5) as feed_name;