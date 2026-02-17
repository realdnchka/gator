-- name: GetPosts :many
SELECT 
    posts.*, 
    ff.user_id 
FROM posts 
JOIN feed_follows AS ff ON posts.feed_id = ff.feed_id 
WHERE ff.user_id = $1
ORDER BY published_at DESC;

-- name: GetPostsByLimit :many
SELECT 
    posts.*, 
    ff.user_id 
FROM posts 
JOIN feed_follows AS ff ON posts.feed_id = ff.feed_id 
WHERE ff.user_id = $1
ORDER BY published_at DESC
LIMIT $2;