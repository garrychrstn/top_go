-- name: UserGetByUsername :one
select * from users where username = $1;

-- name: UserCreate :one
insert into users (username, password) values ($1, $2) returning *;

-- name: UserList :many
select * from users;

-- name: UserPasswordReset :one
update users set password = $1 where id = $2 returning *;
