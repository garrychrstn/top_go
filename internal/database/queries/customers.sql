-- name: CustomerGetByName :one
SELECT * FROM CUSTOMERS WHERE name ILIKE $1;

-- name: CustomerGetByID :one
SELECT * FROM CUSTOMERS WHERE id = $1;

-- name: CustomerList :many
SELECT * FROM CUSTOMERS;

-- name: CustomerCreate :one
INSERT INTO CUSTOMERS (name, phone_number, address) VALUES ($1, $2, $3) returning *;

-- name: CustomerUpdate :one
UPDATE CUSTOMERS SET name = $1, phone_number = $2, address = $3 WHERE id = $4 returning *;

-- name: CustomerDelete :exec
DELETE FROM CUSTOMERS WHERE id = $1;
