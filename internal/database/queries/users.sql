-- name: GetUserByEmail :one
SELECT * FROM USERS WHERE username = $1;

-- name: CreateUser :one
INSERT INTO USERS (username, password) VALUES ($1, $2) returning *;
