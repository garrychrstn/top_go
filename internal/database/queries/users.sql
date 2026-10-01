-- name: UserGetByEmail :one
SELECT * FROM USERS WHERE username = $1;

-- name: UserCreate :one
INSERT INTO USERS (username, password) VALUES ($1, $2) returning *;

-- name: UserList :many
SELECT * FROM USERS;
