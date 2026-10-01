package repository

import (
	"context"

	"github.com/garrychrstn/top-go/internal/repository/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	GetUser(ctx context.Context, username string)
	UpdatePassword(ctx context.Context, id uuid.UUID, password string)
}

type repo struct {
	query *db.Queries
}

func AuthInitRepo(pool *pgxpool.Pool) AuthRepository {
	return &repo{query: db.New(pool)}
}

func (r *repo) GetUser(ctx context.Context, username string) {
}
func (r *repo) UpdatePassword(ctx context.Context, id uuid.UUID, password string) {

}
