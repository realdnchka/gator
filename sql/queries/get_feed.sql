-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE feeds.url = $1 LIMIT 1;