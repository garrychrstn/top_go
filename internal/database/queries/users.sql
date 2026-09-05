-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC, id DESC LIMIT 100;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (name, email) VALUES ($1, $2) RETURNING *;
