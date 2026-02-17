-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows AS ff USING feeds AS f WHERE f.id = ff.feed_id AND f.url = $1 AND f.user_id = $2;