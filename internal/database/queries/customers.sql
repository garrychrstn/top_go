-- name: CustomerGetByName :one
SELECT * FROM CUSTOMERS WHERE name ILIKE $1;

-- name: CustomerCreate :one
INSERT INTO CUSTOMERS (name, phone_number, address) VALUES ($1, $2, $3) returning *;
