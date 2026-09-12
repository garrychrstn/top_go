# go-v1

Chi + sqlc boilerplate. Layered API: `routes -> middleware -> handler -> repository`.
The service layer is intentionally dropped for v1 — handlers call repositories directly.

## Project structure

```
.
├── cmd
│   └── main.go              # wiring: env -> db -> repo -> handler -> router -> server
├── sqlc.yaml                # sqlc codegen config
├── internal
│   ├── database
│   │   ├── db.go            # pgx pool + migration runner
│   │   ├── migrations       # schema, one SQL statement per file (doubles as sqlc schema)
│   │   └── queries          # *.sql query files, consumed by sqlc
│   ├── handler
│   │   ├── middleware       # logging, recover, CORS, request id, auth
│   │   ├── health.go
│   │   ├── router.go        # middleware stack + route mounting
│   │   └── users.go         # sample CRUD handler -> repository
│   ├── repository
│   │   ├── db               # sqlc-generated code (commit it, regenerate via `sqlc generate`)
│   │   ├── errors.go        # sentinel errors + pg error mapping
│   │   └── users.go         # UserRepository interface + pgx implementation
│   ├── types                # shared DTOs / request-response models
│   └── util                 # env loading, JSON writers, decode helpers
```

## Prerequisites

- Go 1.26+
- Postgres 13+ (for `gen_random_uuid()`-free demo; the sample uses `bigserial`)
- [sqlc](https://sqlc.dev) (v1.31+, config is `version: "2"`)

## Quickstart

```sh
# 1. Configure
cp .env.example .env          # set DATABASE_URL to a reachable Postgres

# 2. (Re)generate sqlc code from internal/database/queries
sqlc generate

# 3. Run — connects to DB, applies pending migrations, starts the server
go run ./cmd
```

Migrations run automatically on startup (tracked in a `schema_migrations` table).
To add a table: create `internal/database/migrations/000N_name.sql` (one statement
per file — it is parsed both by sqlc and the migration runner), add queries to
`internal/database/queries`, then `sqlc generate`.

## Endpoints

| Method | Path          | Description          |
| ------ | ------------- | -------------------- |
| GET    | /health       | Liveness check       |
| GET    | /users        | List users           |
| POST   | /users        | Create a user        |
| GET    | /users/{id}   | Get a user by id     |

```sh
curl localhost:8080/health
curl -X POST localhost:8080/users -H 'Content-Type: application/json' \
  -d '{"name":"Ada","email":"ada@example.com"}'
curl localhost:8080/users/1
```

## Notes

- The `users` entity is a sample that demonstrates the full chain end-to-end.
  Replace it (table, queries, repository, handler) with your real domain.
- Handlers implement the `handler.Registrar` interface (`Register(chi.Router)`);
  `NewRouter` mounts any number of them via `handlers := []handler.Registrar{...}`.
- `internal/repository/db` is sqlc output — regenerate after editing queries,
  and commit it like any other code.
- `internal/controller` is not part of the v1 flow (handlers call repositories
  directly) and can be deleted.
