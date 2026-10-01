-- name: CreateItem :one
INSERT INTO items (name, price) VALUES ($1, $2) RETURNING *;

-- name: ListItems :many
SELECT * FROM items;

-- name: UpdateItem :one
UPDATE items SET name = $1, price = $2 WHERE id = $3 RETURNING *;
