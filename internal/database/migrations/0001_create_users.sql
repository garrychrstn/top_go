-- Sample table demonstrating the full handler -> repository -> sqlc chain.
-- One SQL statement per file (see internal/database/db.go).
CREATE TABLE users (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    username    TEXT NOT NULL,
    password    TEXT NOT NULL,
    email      TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
);
