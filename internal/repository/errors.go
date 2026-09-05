package repository

import "errors"

// Sentinel errors returned by repositories so handlers can map them to
// HTTP status codes without depending on Postgres/pgx specifics.
var (
	ErrNotFound   = errors.New("not found")
	ErrEmailTaken = errors.New("email already taken")
)
