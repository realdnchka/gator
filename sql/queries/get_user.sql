-- name: GetUserByName :one
SELECT * FROM users WHERE users.name = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE users.id = $1 LIMIT 1;