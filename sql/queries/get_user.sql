-- name: GetUser :one
SELECT * FROM users WHERE users.name = $1 LIMIT 1;
