# SKILLS.md — go-v1

Chi + sqlc boilerplate API. Go 1.26, Postgres, pgx/v5. Layered flow: `routes -> middleware -> handler -> repository` (no service layer — handlers call repos directly).

## Layout

| Path | Purpose |
|------|---------|
| `cmd/main.go` | Wiring: env → DB → repo → handler → router → server (graceful shutdown) |
| `internal/database` | `db.go` (pgx pool + migration runner), `migrations/` (schema, 1 statement/file), `queries/` (sqlc input) |
| `internal/repository` | `db/` = sqlc-generated code (committed); `*_repo` interfaces + impls, sentinel errors in `errors.go` |
| `internal/handler` | `router.go` (middleware stack), `middleware/` (recover, request id, slog logging, CORS, bearer auth), per-resource handlers |
| `internal/types` | Shared DTOs (request/response structs) |
| `internal/util` | `LoadDotEnv`, JSON write/error/decode helpers |
| `sqlc.yaml` | Codegen config (v2, pgx/v5, timestamptz → `time.Time` override) |
| `.env.example` | `DATABASE_URL`, `PORT`, `ENV`, `CORS_ALLOWED_ORIGINS`, `AUTH_BEARER_TOKEN` |

`internal/controller/` is dead weight from an earlier scaffold — not part of the v1 flow.

## Commands

```sh
sqlc generate    # regenerate internal/repository/db after editing queries/
go run ./cmd     # connects to DB, auto-applies pending migrations, serves
go build ./...   # or go vet ./... / gofmt
```

## Conventions

- **Adding a feature**: migration file (one SQL statement) → query in `internal/database/queries` → `sqlc generate` → repository interface + impl (map pg errors to sentinels) → handler with `Register(r chi.Router)` → mount via `handler.NewRouter`.
- Migrations auto-apply at startup, tracked in `schema_migrations`; the migrations dir doubles as sqlc's schema source, so keep files single-statement.
- Repositories return `types.*` DTOs, never raw sqlc models; handlers call repos directly, no service layer.
- Errors: sentinel errors (`ErrNotFound`, `ErrEmailTaken`…) map to 404/409 in handlers; everything else logs via slog and returns generic 500.
- API errors are always `{"error": "..."}`; success responses are bare JSON.